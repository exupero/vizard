import d3 from 'd3';
import {svgAsDataUri} from 'save-svg-as-png';

window.d3 = d3;

d3.select('#save').on('click', function() {
  svgAsDataUri(d3.select('main svg').node(), {}, function(uri) {
    var r = new XMLHttpRequest();
    r.open('POST', '/svg', true);
    r.send(uri);
  });
});
