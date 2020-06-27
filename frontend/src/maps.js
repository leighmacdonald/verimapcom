"use strict";
var __spreadArrays = (this && this.__spreadArrays) || function () {
    for (var s = 0, i = 0, il = arguments.length; i < il; i++) s += arguments[i].length;
    for (var r = Array(s), k = 0, i = 0; i < il; i++)
        for (var a = arguments[i], j = 0, jl = a.length; j < jl; j++, k++)
            r[k] = a[j];
    return r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.map_create_example = exports.map_view_mission = exports.map_select_coord = exports.map_create_connect = exports.get_map = exports.get_maps = void 0;
require("ol/ol.css");
var ol_1 = require("ol");
var Tile_1 = require("ol/layer/Tile");
var proj_1 = require("ol/proj");
var TileWMS_1 = require("ol/source/TileWMS");
var BingMaps_1 = require("ol/source/BingMaps");
var XYZ_1 = require("ol/source/XYZ");
var OSM_1 = require("ol/source/OSM");
var Extent_1 = require("ol/interaction/Extent");
var BASE_URL = "http://maps.verimap.com/geoserver";
var DEFAULT_SRS = "EPSG:4326";
var DEFAULT_FMT = "image/png";
function get_maps(callable) {
    var url = 'https://cms.verimap.com/examples?public=true';
    fetch(url)
        .then(function (res) { return res.json(); })
        .then(callable)
        .catch(function (err) {
        throw err;
    });
}
exports.get_maps = get_maps;
function get_map(example_id, callable) {
    var url = "https://cms.verimap.com/examples/" + example_id + "?public=true";
    fetch(url)
        .then(function (res) { return res.json(); })
        .then(callable)
        .catch(function (err) {
        throw err;
    });
}
exports.get_map = get_map;
function map_create_connect(target) {
    var osm = new Tile_1.default({
        source: new OSM_1.default({ 'attributions': null })
    });
    var view = new ol_1.View({
        center: proj_1.fromLonLat([-114.366746, 51.097423]),
        zoom: 14,
    });
    return new ol_1.Map({
        layers: [osm],
        target: target,
        view: view
    });
}
exports.map_create_connect = map_create_connect;
function create_vector_layer(layer, srs, fmt) {
    if (srs === void 0) { srs = DEFAULT_SRS; }
    if (fmt === void 0) { fmt = DEFAULT_FMT; }
    return new Tile_1.default({
        visible: true,
        source: new TileWMS_1.default({
            url: BASE_URL + "/demows/wms",
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
function map_select_coord(target_id, on_update) {
    var view = new ol_1.View({
        center: proj_1.fromLonLat([-114.07, 51.05]),
        zoom: 10,
        maxZoom: 16
    });
    var esri = new Tile_1.default({
        source: new XYZ_1.default({
            url: 'https://server.arcgisonline.com/ArcGIS/rest/services/' +
                'World_Imagery/MapServer/tile/{z}/{y}/{x}'
        })
    });
    var layers = {
        'esri': esri
    };
    var map = new ol_1.Map({
        layers: [layers['esri']],
        target: target_id,
        view: view
    });
    var extent = new Extent_1.default();
    map.addInteraction(extent);
    extent.setActive(false);
    //Enable interaction by holding shift
    window.addEventListener('keydown', function (event) {
        if (event.keyCode === 16) {
            extent.setActive(true);
        }
    });
    window.addEventListener('keyup', function (event) {
        if (event.keyCode === 16) {
            extent.setActive(false);
        }
        on_update(extent, map, view);
    });
    return {
        'extent': extent,
        'map': map,
        'layers': layers,
        'view': view,
    };
}
exports.map_select_coord = map_select_coord;
function map_view_mission(target_id, mission_id, lat, lon, srs) {
    if (srs === void 0) { srs = DEFAULT_SRS; }
    var projection = proj_1.get(srs);
    var vector_layers = [];
    var view = new ol_1.View({
        center: proj_1.fromLonLat([lon, lat]),
        zoom: 10,
        maxZoom: 16
    });
    var esri = new Tile_1.default({
        source: new XYZ_1.default({
            url: 'https://server.arcgisonline.com/ArcGIS/rest/services/' +
                'World_Imagery/MapServer/tile/{z}/{y}/{x}'
        })
    });
    var layers = {
        'esri': esri,
        'bing': new Tile_1.default({
            source: new BingMaps_1.default({
                imagerySet: 'Aerial',
                key: "AgmK4-v3LIF5w4UD0u-_y2Sw393klUG9_mXoENnRVC1XGTO393VWfi9yv5uceXq7"
            })
        })
    };
    var map = new ol_1.Map({
        layers: __spreadArrays([layers['esri']], vector_layers),
        target: target_id,
        view: view
    });
    return {
        'map': map,
        'layers': layers,
        'view': view,
    };
}
exports.map_view_mission = map_view_mission;
function map_create_example(example, srs) {
    if (srs === void 0) { srs = DEFAULT_SRS; }
    var projection = proj_1.get(srs);
    var target = "map";
    var vector_layers = [];
    if (example['vector_layers']) {
        example['vector_layers'].split("\n").forEach(function (layer_name) {
            vector_layers.push(create_vector_layer(layer_name));
        });
    }
    var zoom_min = 0;
    var zoom_max = 16;
    if (example['zoom_min'] && example['zoom_min'] >= zoom_min) {
        zoom_min = example['zoom_min'];
    }
    if (example['zoom_max'] && example['zoom_max'] > 0) {
        zoom_max = example['zoom_max'];
    }
    var view = new ol_1.View({
        center: proj_1.fromLonLat([example.longitude, example.latitude]),
        zoom: example.zoom,
        maxZoom: zoom_max
    });
    console.log("min: " + zoom_min + " max: " + zoom_max);
    var esri = new Tile_1.default({
        source: new XYZ_1.default({
            // attributions: null,
            url: 'https://server.arcgisonline.com/ArcGIS/rest/services/' +
                'World_Imagery/MapServer/tile/{z}/{y}/{x}'
        })
    });
    var layers = {
        'esri': esri,
        'bing': new Tile_1.default({
            source: new BingMaps_1.default({
                imagerySet: 'Aerial',
                key: "AgmK4-v3LIF5w4UD0u-_y2Sw393klUG9_mXoENnRVC1XGTO393VWfi9yv5uceXq7"
            })
        }),
        'example': new Tile_1.default({
            minZoom: zoom_min,
            maxZoom: zoom_max,
            source: new TileWMS_1.default({
                url: BASE_URL + "/wms",
                params: { 'LAYERS': example['layer'], 'TILED': true },
                serverType: 'geoserver',
                projection: projection,
                attributions: [
                    'Tiles © <a href="https://verimap.com">Verimap Plus Inc.</a>'
                ],
                transition: 0
            })
        })
    };
    var map = new ol_1.Map({
        layers: __spreadArrays([layers['esri'], layers['example']], vector_layers),
        target: target,
        view: view
    });
    return {
        'map': map,
        'layers': layers,
        'view': view,
    };
}
exports.map_create_example = map_create_example;
