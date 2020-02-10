// Place your application-specific JavaScript functions and classes here
// This file is automatically included by javascript_include_tag :defaults
String.prototype.capitalize = function() {
        return this.replace( /(^|\s)([a-z])/g , function(m,p1,p2){ return p1+p2.toUpperCase(); } );
};

$(function() {
        var $mybuttons = $(".mybutton");
        if ($mybuttons.length > 0) {
                $mybuttons.button();
        };

        var $accordion = $(".accordion");
        if ($accordion.length > 0) {
                $accordion.accordion( {
                        header: ".header",
                        heightStyle: "fill"
                } );
        };

        $(".button-to").one("click", function() {
                var side = $(this).children().children("input[value=\"Accept\"]").attr("class");
                var action = $(this).attr("action");
        });	

        formatMyForms();
        formatTables();
        scrollGameLog();
        scrollChatBox();
        InitMessageLog();
        var $tips = $("div[data-tip]");
        if ($tips.length > 0) {
                $tips.tooltip({items: "div[data-tip]", content: function(){ return $(this).attr("data-tip")}});
        }
});

function formatMyForms() {
        var $myForms = $(".myForm");
        if ($myForms.length > 0 ) {
                $myForms.each(function() {
                        var wrapped = $(this)
                        var $action = $(this).attr("action");
                var $json = $(this).attr("data-json");
                var $dtype = "html";
                if ($json) {
                        $dtype = "json"
                }
                $(this).ajaxForm({
                        url: $action,
                        dataType: $dtype,
                        success: function(data) {
                                if ($json) {
                                        MyAjax.JSONUpdate(data);
                                } else {
                                        MyAjax.Update(data);
                                }
                        },
                        type: "PUT",
                        error: function(data) {
                                alert("Error");
                                window.location.href = target;
                        }
                });
                });
        }
}

function getPath(obj) {
        var $rawPath = obj.data("path");
        if ($rawPath == undefined) {
                return 'none';
        }

        if ($rawPath.length == 0) {
                return window.location.href;
        }

        if ($rawPath[0] == '/') {
                return $rawPath;
        } else {
                return window.location.href + '/' + $rawPath;
        }
}

function getColumns(obj) {
        var $col = [];
        obj.find("th").each(function() {
                var $name = $(this).attr("abbr");
                $col.push({ data: $name });
        });
        return $col
}

function formatTables() {
        var $strippedDataTables = $('table.strippedDataTable');
        if ($strippedDataTables.length > 0) {
                //$strippedDataTables.each(function() {
                $strippedDataTables.dataTable( {
                        "bJQueryUI": true,
                        "sPaginationType": "full_numbers",
                        "bPaginate": false,
                        "bLengthChange": false,
                        "bFilter": false,
                        "bSort": false,
                        "bInfo": false,
                        "bAutoWidth": false,
                        "oLanguage": {
                                "sEmptyTable": "No data available in table"
                        }
                });
        };

        var $myDataTable = $("table.myDataTable");
        $myDataTable.each(function() {
                var $cols = getColumns($(this));
                var $opts = { processing: true,
                        serverSide: false,
                searching: false,
                ordering: false,
                info: true,
                columns: $cols,
                pagingType: "full_numbers",
                pageLength: 15,
                lengthMenu: [ [15, 30], [15, 30] ]
                };

                var $paging = $(this).data("paging");
                if ( $paging == false ) {
                        $opts.paging = false;
                        $opts.info = false;
                };

                var $path = getPath($(this));
                if ( $path != 'none' ) {
                        $opts.ajax = { url: $path, type: "POST" };
                        $opts.serverSide = true;
                };

                var $t = $(this).dataTable($opts);

                $t.on('draw.dt', function() {
                        $('.mybutton').button();
                });
        });
        };

        function scrollChatBox() {
                var log_content = $("#chatbox .content");
                if(log_content[0]) {
                        log_content.scrollTop(log_content[0].scrollHeight);		
                };
        }

        function scrollGameLog() {
                var log_content = $("#gamelog .content");
                if(log_content[0]) {
                        log_content.scrollTop(log_content[0].scrollHeight);		
                };
        }

        function InitMessageLog() {
                var chatbox = $("#chatbox");
                $(".mybutton.post-new-message").one("click", function() {
                        var message = chatbox.find("textarea#message").val();
                        if (message == "") {
                                InitMessageLog();
                                return;
                        }
                        var id = chatbox.attr("class");
                        var target = window.location.href;
                        var targetSplit = target.split("/");
                        var popped = targetSplit.pop();
                        if (popped == "show") {
                                target = targetSplit.join("/");
                        }
                        var creatorid = chatbox.find("div#creatorid").attr("class");
                        var data = {"action" : "post-new-message", "message" : message, "creatorid" : creatorid};
                        $.ajax({
                                url: target + "/addmessage",
                                data: data,
                                dataType: "html",
                                success: function(data) {
                                        $("#chatbox .content").append(data);
                                        $("#chatbox textarea#message").val("");
                                        scrollChatBox();
                                        InitMessageLog();
                                },
                                type: "PUT",
                                error: function(data) {
                                        window.location.href = target;
                                }
                        }); 

                });
        };

        $(window).scroll(function () {
                var $dialog = $("#dialog");
                if ($dialog.length > 0 && !$dialog.dialog("option", "modal")) {
                        if ($dialog.dialog( "isOpen" )) {
                                $dialog.dialog({
                                        position: {
                                                my: "right top",
                                        at: "right top",
                                        of: $(window),
                                        within: "#board" 
                                        }
                                });
                        };
                };
        });

        function DisplayDialog(data) {
                var $dialog = $(data).filter("#dialog");
                InitDialog($dialog);
        };

        function InitDialog(d) {
                var $dialog = d;
                $dialog.find("#tabs").tabs({ selected: 0});
                $dialog.find(".mybutton").button().css("width", "125px");
                $dialog.find(".radio").buttonset();
                var $tips = $dialog.find("div[data-tip]");
                $tips.tooltip({items: "div[data-tip]", content: function(){ return $(this).attr("data-tip")}});
                $dialog.find(".mybutton.submit").one("click", function() {
                        $(this).parent(".myForm").submit();
                });

                var $ajaxTarget = $dialog.attr("action");
                var $json = $dialog.find("input[name='json']");
                var $dtype = "html";
                if ($json.length > 0) {
                        $dtype = "json";
                }
                $dialog.find('.myForm').ajaxForm({
                        url: $ajaxTarget,
                        dataType: $dtype,
                        success: function(data) {
                                if ($json.length > 0) {
                                        MyAjax.JSONUpdate(data);
                                } else {
                                        MyAjax.Update(data);
                                }
                        },
                        type: "PUT",
                        error: function(data) {
                                window.location.href = $ajaxTarget;
                        }
                });

                var modalS = $dialog.attr("data-modal");
                if (modalS != undefined && modalS == 'true') {
                        $dialog.dialog({modal:true});
                } else {
                        $dialog.dialog();
                }
                var maxHeightS = $dialog.attr("data-maxHeight");
                if (maxHeightS != undefined) {
                        var maxHeight = parseInt(maxHeightS);
                        $dialog.dialog("option", "maxHeight", maxHeight);
                }
                var minWidthS = $dialog.attr("data-minWidth");
                if (minWidthS != undefined) {
                        var minWidth = parseInt(minWidthS);
                        $dialog.dialog("option", "minWidth", minWidth);
                }
                var $content = $dialog.find(".dialog-content");
                if ($content.length != 0) {
                        var title = $dialog.dialog("widget").find(".ui-dialog-title");
                        var titleLength = title.text().length;
                        var titleWidth = Math.floor(0.6 * titleLength * parseInt(title.css("font-size").split("px")[0])) + 34;
                        var maxContentWidth = 0;
                        $content.each(function() {
                                var width = parseInt($(this).css("width").split("px")[0]) + 27;
                                if (width > maxContentWidth) {
                                        maxContentWidth = width;
                                };
                        });
                        var $width = Math.max(maxContentWidth, titleWidth);
                        $dialog.dialog("option", "width", $width);
                };
                $(window).scroll();
        };

        var MyAjax = { 
                Update: function(data) {
                        this.update(data); 
                },
                update: function(data) {
                        $("#dialog").dialog("destroy");
                        $("#dialog").remove();
                        data = $.parseHTML(data);
                        var $dialog = $(data).filter("#dialog");
                        var $update = $(data).not("#dialog");
                        $update.children().each(function() {
                                var $id = "#" + $(this).attr("id");
                                $($id).replaceWith($(this));
                        });
                        DisplayDialog($dialog);
                        scrollGameLog();
                        InitClickables();
                }
        };

        ////////////////////////////////////////////////////////////////////////////////////////////
        // Hide Action Forms
        ////////////////////////////////////////////////////////////////////////////////////////////
        $(function () {
                $(".actions").hide();
                //        $(".mybutton").button();
                $(".player-action div.action-link")
                .add(".admin-action")
                .bind('click', generalActionHandler);
        });

        function generalActionHandler(event) {
                var action = $(event.target).parents(".admin-action").attr("id");
                $(".actions").hide();
                $(".actions#"+action+"-action").show();
        };

        $(function() {
                var $dialogData = $("#dialog");
                if ($dialogData.children().length > 0) {
                        DisplayDialog($dialogData);
                };

                //$(".clickmap").mapster({
                //        scaleMap: true,
                //        fill: false,
                //        stroke: true,
                //        strokeWidth: 10,
                //        strokeOpacity: 0.75,
                //        strokeColor: 'ffff00',
                //        onClick: mapClickHandler,
                //        mapKey: 'key',
                //        singleSelect: true,
                //});
                InitClickables();
        });

        function mapClickHandler(data) {
                var $ajaxTarget = window.location.href;
                $.ajax({
                        url: $ajaxTarget,
                        data: { "area" : data.key, "action" : "select-area"}, 
                        dataType: "html",
                        success: function(data) {
                                MyAjax.Update(data);
                        },
                        type: "PUT",
                        error: function(data) {
                                alert("Error");
                                window.location.href = target;
                        }
                }); 
        }

        function InitClickables() {
                $(".clickish").off("click").one("click", function() {
                        $(this).removeClass("clickish");
                        var $area = $(this).attr("id");
                        var $ajaxTarget = window.location.href;
                        $.ajax({
                                url: $ajaxTarget,
                                data: { "area" : $area, "action" : "select-area"}, 
                                dataType: "json",
                                success: function(data) {
                                        MyAjax.JSONUpdate(data);
                                },
                                type: "PUT",
                                error: function(data) {
                                        alert("Error: " + $(data));
                                        window.location.href = $ajaxTarget;
                                }
                        }); 
                });
                $(".clickish").css({"border-style": "solid"}).animate({
                        "border-color": "yellow",
                        "border-width": "3px",
                        "border-style": "solid"
                }, 700, "swing");
                $(".clickable").off("click").one("click", function() {
                        $(this).removeClass("clickable");
                        var $area = $(this).attr("id");
                        var $ajaxTarget = window.location.href;
                        $.ajax({
                                url: $ajaxTarget,
                                data: { "area" : $area, "action" : "select-area"}, 
                                dataType: "html",
                                success: function(data) {
                                        MyAjax.Update(data);
                                },
                                type: "PUT",
                                error: function(data) {
                                        alert("Error: " + $(data));
                                        window.location.href = $ajaxTarget;
                                }
                        }); 
                });
                $(".clickable").css({"border-style": "solid"}).animate({
                        "border-color": "yellow",
                        "border-width": "3px",
                        "border-style": "solid"
                }, 700, "swing");
                $(".multiclickable").on("click", function() {
                        var $area = $(this).attr("id");
                        var $ajaxTarget = window.location.href;
                        $.ajax({
                                url: $ajaxTarget,
                                data: { "area" : $area, "action" : "select-area"}, 
                                dataType: "html",
                                success: function(data) {
                                        MyAjax.Update(data);
                                },
                                type: "PUT",
                                error: function(data) {
                                        alert("Error");
                                        window.location.href = target;
                                }
                        }); 
                });
                $(".multiclickable").animate({
                        "border-color": "yellow",
                        "border-width": "5px",
                        "border-style": "solid"
                }, 700, "swing");
        }

        function trim11 (str) {
                str = str.replace(/^\s+/, '');
                for (var i = str.length - 1; i >= 0; i--) {
                        if (/\S/.test(str.charAt(i))) {
                                str = str.substring(0, i + 1);
                                break;
                        }
                }
                return str;
        }
