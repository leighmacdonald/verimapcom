import '../scss/foundation.scss';
import '../scss/app.scss';
import 'foundation-sites/dist/js/plugins/foundation.core';
import 'foundation-sites/dist/js/plugins/foundation.orbit';
import 'foundation-sites/dist/js/plugins/foundation.smoothScroll';
import 'foundation-sites/dist/js/plugins/foundation.sticky';
import 'foundation-sites/dist/js/plugins/foundation.equalizer';
import 'foundation-sites/dist/js/plugins/foundation.tabs';
import 'foundation-sites/dist/js/plugins/foundation.dropdownMenu';
import 'foundation-sites/dist/js/plugins/foundation.util.keyboard';
import 'foundation-sites/dist/js/plugins/foundation.util.box';
import 'foundation-sites/dist/js/plugins/foundation.util.timer';
import 'foundation-sites/dist/js/plugins/foundation.util.imageLoader';
import 'foundation-sites/dist/js/plugins/foundation.util.touch';
import 'foundation-sites/dist/js/plugins/foundation.util.nest';
import 'foundation-sites/dist/js/plugins/foundation.util.mediaQuery';
import 'foundation-sites/dist/js/plugins/foundation.util.triggers';
import 'foundation-sites/dist/js/plugins/foundation.util.motion';
import 'foundation-sites/dist/js/plugins/foundation.responsiveMenu';
import 'foundation-sites/dist/js/plugins/foundation.responsiveToggle';
import 'jquery'
import 'what-input'

import {get_map, get_maps, map_create_connect, map_create_example} from "./maps";

let map_sets = [];

function init_examples() {
    get_maps((resp) => {
        resp.forEach((example) => {
            map_sets.push(map_create_example(example));
        })
    });

    jQuery("#collapsing-tabs").on("change.zf.tabs", () => {
        jQuery(".tab_set").each(() => {
            // We lazily just re-render all the containers on tab change
            map_sets.forEach((ms) => {
                ms.map.updateSize();
            });
        });
    });
}

function main() {
    jQuery(document).foundation();
    const path = window.location.pathname.toLowerCase();
    if (path === "/connect") {
        map_create_connect("map");
        return;
    }
    let match = path.match(/\/example\/(\d+)/);
    if (match) {
        get_map(match[1], map_create_example);
    }
}

document.addEventListener("DOMContentLoaded", main);
