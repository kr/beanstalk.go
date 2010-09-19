/* Modified from godocs.js in the go distribution.
 * Generates a table of contents: looks for h2 and h3 elements and generates
 * links. "Decorates" the element with id=="nav" with this table of contents.
 */
function gen_toc() {
  var navbar = $('#toc');
  if (navbar.length < 1) return;

  var toc_items = $('h2[id], h3[id]').map(function (i, x) {
    var item = $(x.tagName.toLowerCase() == 'h2' ? '<dt>' : '<dd>')
    item.append($('<a>').attr('href', '#'+x.id).text($(x).text()));
    navbar.append(item);
  });
}

$(document).ready(gen_toc);
