define(['knockout-3.2.0', 'jquery', 'jquery-ui/dialog'], function(ko, $, dialog){
        function positionDialog(d) {
                var boardTop = $("#board").offset().top;
                var windowTop = $(window).scrollTop();
                var topOffset = 0;
                if (windowTop < boardTop) {
                        topOffset = boardTop - windowTop;
                }
                d.dialog("option", "position", {
                        my: "left top",
                        at: "left top+" + topOffset,
                        of: window,
                        within: "#board", 
                        collision: "fit",
                }
                );
        }

        $(window).scroll(function () {
                var $dialog = $("#dialog");
                if ($dialog.length > 0 && !$dialog.dialog("option", "modal") && $dialog.dialog("isOpen")) {
                        positionDialog($dialog);
                };
        });


        ko.bindingHandlers.jqDialog = {
                init: function(element, valueAccessor, allBindingsAccessor) {
                        var options = ko.utils.unwrapObservable(valueAccessor());
                        $(element).on("dialogcreate", function(event, ui) {
                                options.initialized(true);
                        });

                        $(element).dialog({ autoOpen: false });

                        //handle disposal (not strictly necessary in this scenario)
                        ko.utils.domNodeDisposal.addDisposeCallback(element, function() {
                                $(element).dialog("destroy");
                        });   
                },
                update: function(element, valueAccessor, allBindingsAccessor) {
                        var options = ko.utils.unwrapObservable(valueAccessor()),
                        initialized = ko.utils.unwrapObservable(options.initialized()) || false;

                        //don't call before initilization
                        if (initialized) {
                                var title = ko.utils.unwrapObservable(options.title()) || 'Default',
                                        open = ko.utils.unwrapObservable(options.open()) || false,
                                             minWidth = ko.utils.unwrapObservable(options.minWidth()) || 200,
                                             $el = $(element);

                                $el.dialog("option", { title: title, minWidth: minWidth });
                                if (open) {
                                        positionDialog($el);
                                        $el.dialog("open");
                                }
                        }  
                }
        };
});
