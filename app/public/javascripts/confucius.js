$(function () {
        $(".mybutton.undo").click(function() {
                $("#undo_turn").submit();
        });
});

////////////////////////////////////////////////////////////////////////////////////////////
// Hide Action Forms
////////////////////////////////////////////////////////////////////////////////////////////
$(function () {
        $(".actions").hide();
        $(".player-action div.action-link")
                .add(".admin-action a")
                .bind('click', generalActionHandler);
});

function generalActionHandler(event) {
        var action = event.target.parentNode.id;
        $(".actions").hide();
        $(".actions#"+action+"-action").show();
};

///////////////////////////////////////////////////////////////////////
// If only one action available, trigger click to show form.
///////////////////////////////////////////////////////////////////////
$(function () {
        if ($(".player-action a").size() == 1) {
                $(".player-action a").triggerHandler('click');
        };
});

$(function () {
        for (var i = 0; i <= 9; i++) {
                emperorDialog($("div.emperorcard"+i), $("img.emperorcard"+i));
        };
});

function emperorDialog(div, img) {
        if (div.length == 1) {
                div.dialog({ 
                        autoOpen: false,
                        title: div.attr("title"),
                });
        }
        if (img.length == 1) {
                img.hover(
                                function () {
                                        div.dialog("open");
                                },
                                function () {
                                        div.dialog("close");
                                }
                         );
        };
}

$(function () {
        for (var i = 1; i <= 2; i++) {
                playerAidDialog($("div.playeraid"+i), $("img.playeraid"+i));
        };
});

function playerAidDialog(div, img) {
        if (div.length == 1) {
                div.dialog({ 
                        autoOpen: false,
                        title: div.attr("title"),
                        modal: true,
                        width: 560,
                });
        }
        if (img.length == 1) {
                img.click(
                                function () {
                                        div.dialog("open");
                                }
                         );
        };
}

function say(text) {
        $('#console').append('<div>'+text+'</div>');
}

// Form Support for Emperor's Reward Action
$(function () {
        $("select#reward-card")
        .bind("change", function(event) {
                var card = $("select#reward-card").val();
                $(".reward-directions").add(".reward-form").hide();
                $(".reward-directions#reward-card-"+card).add(".reward-form#reward-card-"+card).show();
        })
.triggerHandler("change");
});

// Form Support for Petition Emperor Action
$(function () {
        $("select#petition-gift")
        .bind("change", function(event) {
                var card = $("select#petition-gift").val();
                $(".petition-directions").add(".petition-form").hide();
                $("#petition-directions-"+card).add("#petition-form-"+card).show();
        })
.triggerHandler("change");
});

$(function() {
        $(".icons #options")
        .mouseenter(function() {
                $(this).addClass('ui-state-hover');
        })
.mouseleave(function() {
        $(this).removeClass("ui-state-hover");
})
.click(function() {
        $(this).addClass('ui-state-active');
        $(".actions").hide();
        $(".actions#options-action").show();
});
});
