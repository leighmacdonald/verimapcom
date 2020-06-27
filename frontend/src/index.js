"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
require("./scss/foundation.scss");
require("./scss/app.scss");
require("foundation-sites/dist/js/plugins/foundation.core");
require("foundation-sites/dist/js/plugins/foundation.orbit");
require("foundation-sites/dist/js/plugins/foundation.smoothScroll");
require("foundation-sites/dist/js/plugins/foundation.sticky");
require("foundation-sites/dist/js/plugins/foundation.equalizer");
require("foundation-sites/dist/js/plugins/foundation.tabs");
require("foundation-sites/dist/js/plugins/foundation.dropdownMenu");
require("foundation-sites/dist/js/plugins/foundation.util.keyboard");
require("foundation-sites/dist/js/plugins/foundation.util.box");
require("foundation-sites/dist/js/plugins/foundation.util.timer");
require("foundation-sites/dist/js/plugins/foundation.util.imageLoader");
require("foundation-sites/dist/js/plugins/foundation.util.touch");
require("foundation-sites/dist/js/plugins/foundation.util.nest");
require("foundation-sites/dist/js/plugins/foundation.util.mediaQuery");
require("foundation-sites/dist/js/plugins/foundation.util.triggers");
require("foundation-sites/dist/js/plugins/foundation.util.motion");
require("foundation-sites/dist/js/plugins/foundation.responsiveMenu");
require("foundation-sites/dist/js/plugins/foundation.responsiveToggle");
var jquery_1 = require("jquery");
require("what-input");
globalThis.jQuery = jquery_1.default;
var maps_1 = require("./maps");
var proj_1 = require("ol/proj");
var mission_1 = require("./mission");
var map_sets = [];
function init_examples() {
    maps_1.get_maps(function (resp) {
        resp.forEach(function (example) {
            map_sets.push(maps_1.map_create_example(example));
        });
    });
    jquery_1.default("#collapsing-tabs").on("change.zf.tabs", function () {
        jquery_1.default(".tab_set").each(function () {
            // We lazily just re-render all the containers on tab change
            map_sets.forEach(function (ms) {
                ms.map.updateSize();
            });
        });
    });
}
function page_missions_create() {
    var lat_ul = document.getElementById("lat_ul");
    var lon_ul = document.getElementById("lon_ul");
    var lat_lr = document.getElementById("lat_lr");
    var lon_lr = document.getElementById("lon_lr");
    maps_1.map_select_coord("map", function (extent, map, view) {
        var lonlat_ul = proj_1.transform([extent.extent_[0], extent.extent_[1]], 'EPSG:3857', 'EPSG:4326');
        lon_ul.value = lonlat_ul[0];
        lat_ul.value = lonlat_ul[1];
        var lonlat_lr = proj_1.transform([extent.extent_[2], extent.extent_[3]], 'EPSG:3857', 'EPSG:4326');
        lon_lr.value = lonlat_lr[0];
        lat_lr.value = lonlat_lr[1];
    });
}
function main() {
    jquery_1.default(document).foundation();
    var path = window.location.pathname.toLowerCase();
    if (path === "/connect") {
        maps_1.map_create_connect("map");
        return;
    }
    if (path === "/missions/create") {
        page_missions_create();
    }
    var m = path.match(/\/mission\/(\d+)/);
    if (m) {
        mission_1.page_mission(m[1]);
    }
    var match = path.match(/\/example\/(\d+)/);
    if (match) {
        maps_1.get_map(match[1], maps_1.map_create_example);
    }
}
document.addEventListener("DOMContentLoaded", main);
