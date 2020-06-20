import 'ol/ol.css';
import {Map, View} from 'ol';
import TileLayer from 'ol/layer/Tile';
import {fromLonLat, get as getProjection} from 'ol/proj';
import TileWMS from "ol/source/TileWMS";
import BingMaps from "ol/source/BingMaps";
import XYZ from "ol/source/XYZ";
import OSM from "ol/source/OSM";

const BASE_URL = "http://maps.verimap.com/geoserver";
const DEFAULT_SRS = "EPSG:4326";
const DEFAULT_FMT = "image/png";

export function get_maps(callable) {
    let url = 'https://cms.verimap.com/examples?public=true';
    fetch(url)
        .then(res => res.json())
        .then(callable)
        .catch(err => {
            throw err
        });
}

export function get_map(example_id, callable) {
    let url = `https://cms.verimap.com/examples/${example_id}?public=true`;
    fetch(url)
        .then(res => res.json())
        .then(callable)
        .catch(err => {
            throw err
        });
}

export function map_create_connect(target) {
    let osm = new TileLayer({
        source: new OSM({'attributions': null})
    });
    let view = new View({
        center: fromLonLat([-114.366746, 51.097423]),
        zoom: 14,
    });
    return new Map({
        layers: [osm],
        target: target,
        view: view
    });
}

function create_vector_layer(layer, srs = DEFAULT_SRS, fmt = DEFAULT_FMT) {
    return new TileLayer({
        visible: true,
        source: new TileWMS({
            url: `${BASE_URL}/demows/wms`,
            params: {
                'FORMAT': fmt,
                'VERSION': '1.3.0',
                tiled: true,
                "LAYERS": layer,
                "exceptions": 'application/vnd.ogc.se_inimage',
                tilesOrigin: -117.44047995726982 + "," + 59.54844821883142
            }
        })
    });
    // let vectorSource = new VectorSource({
    //     format: new GeoJSON(),
    //     url: function (extent) {
    //         const u = 'https://maps.verimap.com/geoserver/wfs?service=WFS&' +
    //             `version=1.1.0&request=GetFeature&typename=${layer}&` +
    //             `outputFormat=application/json&srsname=${srs}&` +
    //             `bbox=${extent.join(',')},${srs}`;
    //         console.log(u);
    //         return u;
    //     },
    //     strategy: bbox
    // });
    // return new VectorLayer({source: vectorSource});
}

export function map_create_example(example, srs = DEFAULT_SRS) {
    const projection = getProjection(srs);
    const target = `map`;
    const vector_layers = [];
    if (example['vector_layers']) {
        example['vector_layers'].split("\n").forEach((layer_name) => {
            vector_layers.push(create_vector_layer(layer_name))
        });
    }
    let zoom_min = 0;
    let zoom_max = 16;
    if (example['zoom_min'] && example['zoom_min'] >= zoom_min) {
        zoom_min = example['zoom_min'];
    }
    if (example['zoom_max'] && example['zoom_max'] > 0) {
        zoom_max = example['zoom_max'];
    }
    let view = new View({
        center: fromLonLat([example.longitude, example.latitude]),
        zoom: example.zoom,
        maxZoom: zoom_max
    });

    console.log(`min: ${zoom_min} max: ${zoom_max}`);
    const esri = new TileLayer({
        source: new XYZ({
            // attributions: null,
            url: 'https://server.arcgisonline.com/ArcGIS/rest/services/' +
                'World_Imagery/MapServer/tile/{z}/{y}/{x}'
        })
    });
    const layers = {
        'esri': esri,
        'bing': new TileLayer({
            source: new BingMaps({
                imagerySet: 'Aerial',
                key: `AgmK4-v3LIF5w4UD0u-_y2Sw393klUG9_mXoENnRVC1XGTO393VWfi9yv5uceXq7`
            })
        }),
        'example': new TileLayer({
            minZoom: zoom_min,
            maxZoom: zoom_max,
            source: new TileWMS({
                url: `${BASE_URL}/wms`,
                params: {'LAYERS': example['layer'], 'TILED': true},
                serverType: 'geoserver',
                projection: projection,
                attributions: [
                    'Tiles © <a href="https://verimap.com">Verimap Plus Inc.</a>'
                ],
                transition: 0
            })
        })
    };
    const map = new Map({
        layers: [layers['esri'], layers['example'], ...vector_layers],
        target: target,
        view: view
    });
    return {
        'map': map,
        'layers': layers,
        'view': view,
    };
}
