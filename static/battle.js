var _href = $("#edit-button").attr("href");

function onChange() {
  $(function () {
    var embed = "";
    $(".playButton").click(function () {
      var button = $(this);
      embed = button.data("embed");
      var toembed = button.closest(".embedded-track");
      toembed.html(embed);
      expandWaveform();
      toembed.find("iframe").css("width", "100%");
    });
  });
  if ($("zing-grid").hasClass("zing-vote")) {
    $("body > div.container > zing-grid > zg-head > zg-row > zg-head-cell:nth-child(2)").addClass("zing-vote-header");
  }
  
  $(".form-ajax").unbind("submit")

  $(".form-ajax").submit(function(t) {
    t.preventDefault();
    var e = $(this).attr("action"),
        o = $(this).attr("method"),
        n = $(this).serialize(),
        r = $(this).find(".material-icons");
    i = $(this).find(".material-icons"), row = $(this).closest("zg-row").attr("aria-rowindex"), col = $(this).closest("zg-cell").attr("aria-colindex"), $.ajax({
        url: e,
        type: o,
        data: n,
        success: function(t) {
            t.Redirect ? window.location.replace(t.RedirectPath) : (M.toast({
                html: t.ToastHTML,
                classes: t.ToastClass,
            }), "liked" == t.ToastQuery && i.attr("style", "color: #ff5800"), "unliked" == t.ToastQuery && i.attr("style", ""), "successvote" == t.ToastQuery && ($(".votes-remaining").html(parseInt($(".votes-remaining").html(), 10) - 1), r.attr("style", "color: #ff5800")), "successdelvote" == t.ToastQuery && ($(".votes-remaining").html(parseInt($(".votes-remaining").html(), 10) + 1), r.attr("style", "")))
        }
    })
  })
  $(".tooltipped").tooltip()
  expandWaveform()
}

function loveSort(t, e) {
    return console.log(t), "#ff5800" == t.like_colour ? -1 : "#ff5800" == e.like_colour ? 1 : t.like_colour > e.like_colour ? 1 : -1
}

function voteSort(t, e) {
    return console.log(t), "#ff5800" == t.vote_colour ? -1 : "#ff5800" == e.vote_colour ? 1 : t.vote_colour > e.vote_colour ? 1 : -1
}
$("#edit-button").attr("href", _href + "timezone/" + Intl.DateTimeFormat().resolvedOptions().timeZone);