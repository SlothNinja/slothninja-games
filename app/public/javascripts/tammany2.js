$(function () {
        $(".clickmap").mapster({
                scaleMap: true,
                fill: false,
                stroke: true,
                strokeWidth: 10,
                strokeOpacity: 0.75,
                strokeColor: 'ffff00',
                onClick: mapClickHandler,
                mapKey: 'key',
                singleSelect: true,
        });

        $("th.slander").addClass("ui-state-default");
        $(".increment-button-0").button({ icons : { primary: 'ui-icon-carat-1-n' }, text : false }).click(function(){
                var target1 = $(this).parent().find(".text");
                var target2 = $(this).parent().find("input");
                var count = parseInt(target1.html(), 10);
                var max = parseInt(/max-(\d\d*)/.exec($(this).attr("class"))[1], 10);
                var row = $("tr#bid-0");
                if (count < max) {
                        count += 1;
                        target1.html(count);
                        target2.val(count);
                        var total = 0;
                        row.find(".chip .text").add(row.find(".boss .text")).each(function() {
                                total += parseInt($(this).html(), 10);
                        });
                        row.find(".total").html(total);
                }
        });
        $(".decrement-button-0").button({ icons : { primary: 'ui-icon-carat-1-s' }, text : false }).click(function(){
                var target1 = $(this).parent().find(".text");
                var target2 = $(this).parent().find("input");
                var count = parseInt(target1.html(), 10);
                var row = $("tr#bid-0");
                if (count > 0) {
                        count -= 1;
                        target1.html(count);
                        target2.val(count);
                        var total = 0;
                        row.find(".chip .text").add(row.find(".boss .text")).each(function() {
                                total += parseInt($(this).html(), 10);
                        });
                        row.find(".total").html(total);
                }
        });
        $(".increment-button-1").button({ icons : { primary: 'ui-icon-carat-1-n' }, text : false }).click(function(){
                var target1 = $(this).parent().find(".text");
                var target2 = $(this).parent().find("input");
                var count = parseInt(target1.html(), 10);
                var max = parseInt(/max-(\d\d*)/.exec($(this).attr("class"))[1], 10);
                var row = $("tr#bid-1");
                if (count < max) {
                        count += 1;
                        target1.html(count);
                        target2.val(count);
                        var total = 0;
                        row.find(".chip .text").add(row.find(".boss .text")).each(function() {
                                total += parseInt($(this).html(), 10);
                        });
                        row.find(".total").html(total);
                }
        });
        $(".decrement-button-1").button({ icons : { primary: 'ui-icon-carat-1-s' }, text : false }).click(function(){
                var target1 = $(this).parent().find(".text");
                var target2 = $(this).parent().find("input");
                var count = parseInt(target1.html(), 10);
                var row = $("tr#bid-1");
                if (count > 0) {
                        count -= 1;
                        target1.html(count);
                        target2.val(count);
                        var total = 0;
                        row.find(".chip .text").add(row.find(".boss .text")).each(function() {
                                total += parseInt($(this).html(), 10);
                        });
                        row.find(".total").html(total);
                }
        });
});

function OfficeDialogInit(data) {
        var office = $(data).find("div").attr("class").split("-")[1];
        var officeStrings = { "2" : "Deputy Mayor", "3" : "Council President", "4" : "Chief Of Police", "5" : "Precinct Chairman" };
        var title = "Assign " + officeStrings[office];
        var dialog = $("#dialog");

        $(".mybutton").button();
        dialog.dialog({ 
                title: title,
                minHeight: 0,
                close: function() { $(this).dialog("destroy")}
        });
        $("#board").append('<div id="anchor"></div>');

        var widget = dialog.dialog("widget");
        widget.css("top", "0");
        $("#anchor").replaceWith(widget);
        widget.jScroll({top: 40}).trigger("scroll");

        $(".mybutton.assign-office").click(function() {
                var target = window.location.href
                var ajaxTarget = window.location.href;
        var actionSplit = $(this).parent().attr("class").split("-");
        var pid = actionSplit[1];
        $.ajax({
                url: ajaxTarget,
                data: { "action" : "assign-office", "playerID" : pid, "office" : office }, 
                dataType: "html",
                success: function(data) {
                        AjaxUpdate(data);
                        DialogHide(dialog);
                        AssignOfficesInit();
                },
                type: "PUT",
                error: function(data) {
                        alert("Error");
                        window.location.href = target;
                }
        }); 
        });
        $(".mybutton.undo").click(function() {
                $("form#undo-turn").submit();
        });
        $(".mybutton.finish").click(function() {
                $("form#finish-turn").submit();
        });
};

function AssignOfficesInit() {
        $(".assign-office").click(function(e){
                var office = $(this).find("img").attr("class").split("-")[1];
                var ajaxTarget = window.location.href;
                $.ajax({
                        url: ajaxTarget,
                        data: { "office" : office, "action" : "select-office" },
                        dataType: "html",
                        success: DisplayOfficeDialog,
                        type: "PUT",
                        error: function(data) {
                                alert("Error");
                                window.location.href = target;
                        }
                });
        });
};

function DisplayOfficeDialog(data) {
        var dialogContent = $(data).filter("#dialog");
        $("#dialog").dialog("destroy");
        $("#dialog").replaceWith(dialogContent)
                OfficeDialogInit(data);
};
