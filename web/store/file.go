package store

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path"
	"time"
)

type File struct {
	FileID    int
	PersonID  int
	FileName  string
	FileType  string
	FileSize  int64
	Hash256   []byte
	Data      []byte
	CreatedOn time.Time
	UpdatedOn time.Time
}

func (f File) Path(root string) string {
	return path.Join(root, fmt.Sprintf("%d", f.FileID))
}

func NewFile(reader io.Reader, fileName string, personID int) (*File, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read data from reader")
	}
	var f File
	f.CreatedOn = time.Now()
	f.UpdatedOn = f.CreatedOn
	f.PersonID = personID
	f.FileName = fileName
	f.FileSize = int64(len(data))
	var b []byte
	b2 := sha256.Sum256(data)
	for byt := range b2 {
		b = append(b, byte(byt))
	}
	f.Hash256 = b
	ct, err2 := getFileContentType(data)
	if err2 != nil {
		log.Warnf("Failed to detect content type: %s", fileName)
		ct = "application/octet-stream"
	}
	f.FileType = ct
	f.Data = data
	return &f, nil
}

func getFileContentType(body []byte) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := bytes.NewReader(body).Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(body)

	return contentType, nil
}

func FileSave(ctx context.Context, db *pgxpool.Pool, saveDir string, f *File) error {
	if f.FileID > 0 {
		return fileUpdate(ctx, db, saveDir, f)
	}
	return fileInsert(ctx, db, saveDir, f)
}

// Exists checks for the existence of a file path
func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func FileDelete(ctx context.Context, db *pgxpool.Pool, dir string, f *File) error {
	const q = `DELETE FROM file WHERE file_id = $1`
	if _, err := db.Exec(ctx, q, f.FileID); err != nil {
		return errors.Wrapf(err, "Failed to delete file row")
	}
	fp := path.Join(dir, f.FileName)
	if Exists(fp) {
		if err := os.Remove(fp); err != nil {
			return errors.Wrapf(err, "failed to delete file")
		}
	}
	return nil
}

func fileInsert(ctx context.Context, db *pgxpool.Pool, dir string, f *File) error {
	const q = `
		INSERT INTO file (
			person_id, file_name, file_type, file_size, hash_256, created_on, updated_on
		) VALUES (
		    $1, $2, $3, $4, $5, $6, $7
		) RETURNING file_id`
	if err := db.QueryRow(ctx, q, f.PersonID, f.FileName, f.FileType, f.FileSize,
		f.Hash256, f.CreatedOn, f.UpdatedOn).Scan(&f.FileID); err != nil {
		return err
	}
	if err := fileWrite(ctx, db, dir, f); err != nil {
		return errors.Wrapf(err, "Failed to write new file")
	}
	return nil
}

func fileRoot(root string, fileID int) string {
	rounded := int(math.Round(float64(fileID)/1000) * 1000)
	return path.Join(root, fmt.Sprintf("%d", rounded))
}

func fileWrite(ctx context.Context, db *pgxpool.Pool, dir string, f *File) error {
	root := fileRoot(dir, f.FileID)
	if !Exists(root) {
		if err := os.Mkdir(root, 0755); err != nil {
			return err
		}
	}
	fullPath := path.Join(root, fmt.Sprintf("%d", f.FileID))
	fp, err := os.Create(fullPath)
	if err != nil {
		_ = FileDelete(ctx, db, dir, f)
		return err
	}
	defer func() {
		if err := fp.Close(); err != nil {
			log.Errorf("Failed to close file after insert: %v", err)
		}
	}()
	c, err := fp.Write(f.Data)
	if err != nil {
		_ = FileDelete(ctx, db, dir, f)
		return err
	}
	if int64(c) != f.FileSize {
		return errors.Errorf("File size and write size mismatch: %d - %d", f.FileSize, c)
	}
	return nil
}

func fileUpdate(ctx context.Context, db *pgxpool.Pool, dir string, f *File) error {
	const q = `
		UPDATE file
		SET
			person_id = $2, file_name = $3, file_type = $4, file_size = $5, hash_256 = $6,
		    created_on = $7, updated_on = $8
		WHERE
			file_id = $1
	`
	f.UpdatedOn = time.Now()
	if _, err := db.Exec(ctx, q, f.FileID, f.PersonID, f.FileName, f.FileType, f.FileSize,
		f.Hash256, f.CreatedOn, f.UpdatedOn); err != nil {
		return errors.Wrapf(err, "Failed to update file")
	}
	if err := fileWrite(ctx, db, dir, f); err != nil {
		return errors.Wrapf(err, "Failed to write updated file")
	}
	return nil
}

// FileRead will read in the file from the filesystem and store the data
// in the .Data field.
func FileRead(dir string, f *File) error {
	fp := path.Join(fileRoot(dir, f.FileID), fmt.Sprintf("%d", f.FileID))
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	if int64(len(b)) != f.FileSize {
		return errors.Errorf("File size mismatched")
	}
	f.Data = b
	return nil
}

func FileGetAllMission(ctx context.Context, db *pgxpool.Pool, missionID int) ([]*File, error) {
	const q = `
		SELECT 
			f.file_id, f.person_id, f.file_name, f.file_type, f.file_size, f.hash_256, f.created_on, f.updated_on 
		FROM 
		     file f
		LEFT JOIN mission_file mf on f.file_id = mf.file_id
		WHERE
			mf.mission_id = $1`
	var files []*File
	rows, err := db.Query(ctx, q, missionID)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to query mission files")
	}
	defer rows.Close()
	for rows.Next() {
		var f File
		if err := rows.Scan(&f.FileID, &f.PersonID, &f.FileName, &f.FileType,
			&f.FileSize, &f.Hash256, &f.CreatedOn, &f.UpdatedOn); err != nil {
			return nil, errors.Wrapf(err, "Failed to scan file row")
		}
	}
	return files, nil
}
func FileGet(ctx context.Context, db *pgxpool.Pool, fileID int) (*File, error) {
	const q = `
		SELECT 
			file_id,person_id, file_name, file_type, file_size, hash_256, created_on, updated_on 
		FROM 
		     file
		WHERE
			file_id = $1`
	var f File
	if err := db.QueryRow(ctx, q, fileID).Scan(&f.FileID, &f.PersonID, &f.FileName,
		&f.FileType, &f.FileSize, &f.Hash256, &f.CreatedOn, &f.UpdatedOn); err != nil {
		return nil, errors.Wrapf(err, "Failed to query file")
	}
	return &f, nil
}

func FileRegisterDownload(ctx context.Context, db *pgxpool.Pool, personID int, fileID int) error {
	const q = `INSERT INTO file_downloads (file_id, person_id, created_on) VALUES ($1, $2, $3)`
	if _, err := db.Exec(ctx, q, fileID, personID, time.Now()); err != nil {
		return errors.Wrapf(err, "Failed to register download")
	}
	return nil
}

func FileHaveAccess(ctx context.Context, db *pgxpool.Pool, agencyID int, fileID int, personID int) bool {
	const q = `
		SELECT 1 FROM file f 
		LEFT JOIN agency_file af on f.file_id = af.file_id
		WHERE af.agency_id = $1 AND af.file_id = $2 OR f.person_id = $3
	`
	var val int
	if err := db.QueryRow(ctx, q, agencyID, fileID, personID).Scan(&val); err != nil {
		log.Errorf("Error checking for access permissions: %v", err)
		return false
	}
	return val == 1
}

func FileUploadsGetPaged(ctx context.Context, db *pgxpool.Pool, p Person, limit, offset int) ([]*File, error) {
	const q = `SELECT 
		    f.file_id, f.person_id, f.file_name, f.file_type, f.file_size, f.hash_256, f.created_on, f.updated_on
		FROM file f
		LEFT JOIN agency_file af on f.file_id = af.file_id
		WHERE 
		      f.person_id = $1 OR af.agency_id = $2
		LIMIT $3 OFFSET $4`
	var files []*File
	rows, err := db.Query(ctx, q, p.PersonID, p.AgencyID, limit, offset)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query uploads")
	}
	defer rows.Close()
	for rows.Next() {
		var f File
		if err := rows.Scan(&f.FileID, &f.PersonID, &f.FileName, &f.FileType, &f.FileSize,
			&f.Hash256, &f.CreatedOn, &f.UpdatedOn); err != nil {
			return nil, errors.Wrapf(err, "Failed to scan file row")
		}
		files = append(files, &f)
	}

	return files, nil
}

func FileGetPaged(ctx context.Context, db *pgxpool.Pool, p Person, limit, offset int) ([]*File, error) {
	const q = `
		SELECT 
		    f.file_id, f.person_id, f.file_name, f.file_type, f.file_size, f.hash_256, f.created_on, f.updated_on
		FROM file f
		LEFT JOIN agency_file af on f.file_id = af.file_id
		WHERE 
		      af.agency_id = $1 OR f.person_id = $2
		LIMIT $3 OFFSET $4
		`
	var files []*File
	rows, err := db.Query(ctx, q, p.AgencyID, p.PersonID, limit, offset)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query files")
	}
	defer rows.Close()
	for rows.Next() {
		var f File
		if err := rows.Scan(&f.FileID, &f.PersonID, &f.FileName, &f.FileType, &f.FileSize,
			&f.Hash256, &f.CreatedOn, &f.UpdatedOn); err != nil {
			return nil, errors.Wrapf(err, "Failed to scan file row")
		}
		files = append(files, &f)
	}

	return files, nil
}
