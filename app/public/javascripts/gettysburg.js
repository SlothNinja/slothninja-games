String.prototype.ucfirst = function() {
        return this.charAt(0).toUpperCase() + this.substr(1);
}

String.prototype.lcfirst = function() {
        return this.charAt(0).toLowerCase() + this.substr(1);
}

var gettysburg = function() {
        function side(s) {
                if (s == 2) {
                        return "confederate";
                } else {
                        return "union";
                }
        }

        function unitModel(data) {
                var self = this;
                self.id = data.id;
                self.side = data.side;
                self.rank = data.rank;
                self.kind = data.kind;

                self.hit = ko.observable(data.hit);
                self.damage = ko.observable(data.damage);
                self.newDamage = ko.observable(data.newDamage);
                self.retreatTo = ko.observable(data.retreatTo);
                self.retreat = ko.observable(data.retreat);

                self.link = function() {
                        return "/images/gettysburg/" + self.side + "-" + self.rank + "-" + self.kind + ".png";
                };

                self.totalDamage = ko.computed(function() {
                        return self.damage() + self.newDamage();
                });

                self.cubes = ko.computed(function() {
                        var cubes = [];
                        for (var i = 0; i < self.totalDamage(); i++) {
                                cubes.push(new cubeModel({side: self.side, rank: self.rank, kind: self.kind}));
                        }
                        return cubes
                });

                self.idString = function(first, second) {
                        return first + "-" + second;
                }

                self.myUpdate = function(data) {
                        if (!data) {
                                return
                        }
                        if (self.hit() != data.hit) {
                                self.hit(data.hit)
                        }
                        if (self.damage() != data.damage) {
                                self.damage(data.damage)
                        }
                        if (self.newDamage() != data.newDamage) {
                                self.newDamage(data.newDamage)
                        }
                        if (self.retreatTo() != data.retreatTo) {
                                self.retreatTo(data.retreatTo)
                        }
                        if (self.retreat() != data.retreat) {
                                self.retreat(data.retreat)
                        }
                }
        }

        function cubeModel(data) {
                var self = this;
                self.side = data.side;
                self.rank = data.rank;
                self.kind = data.kind;
                self.link = function() {
                        return "/images/gettysburg/" + self.side + "-" + self.rank + "-" + self.kind + "-cube.png";
                };

        }

        function areaModel(data) {
                var self = this;
                self.id = data.id;
                self.key = "area-" + self.id;
                self.hit = ko.observable(data.hits);

                self.commandBlock = ko.observable(data.commandBlock ?  new commandBlockModel(data.commandBlock) : null);
                self.toUnits = function(data) {
                        var units = [];
                        if (!data) {
                                return [];
                        }
                        for (var i = 0; i < data.length; i++) {
                                units.push(ko.observable(new unitModel(data[i])));
                        }
                        return units;
                };
                self.units = ko.observableArray(self.toUnits(data.units));
                self.units.changed = function(data) {
                        var units = self.units();
                        if (units.length == 0 && data == null) {
                                return false;
                        }
                        if (units.length > 0 && data == null) {
                                return true;
                        }
                        if (units.length != data.length) {
                                return true;
                        }
                        for (var i = 0; i < units.length; i++) {
                                var unit = units[i]();
                                if (unit.id != data[i].id) {
                                        return true;
                                }
                                if (unit.damage() != data[i].damage) {
                                        return true;
                                }
                                if (unit.newDamage() != data[i].newDamage) {
                                        return true;
                                }
                        }
                        return false;
                };

                self.units.myUpdate = function(data) {
                        if (self.units.changed(data)) {
                                self.units(self.toUnits(data));
                        }
                }

                self.orderDiscs = ko.observable(data.orderDiscs.data);
                self.orderDiscs.changed = function(data) {
                        if (self.orderDiscs().confederateSingle != data.confederateSingle) {
                                return true;
                        }
                        if (self.orderDiscs().confederateDouble != data.confederateDouble) {
                                return true;
                        }
                        if (self.orderDiscs().unionSingle != data.unionSingle) {
                                return true;
                        }
                        if (self.orderDiscs().unionForcedPass != data.unionForcedPass) {
                                return true;
                        }
                        return false;
                };
                self.discs = ko.computed(function() {
                        return initDiscs(self.orderDiscs());
                });

                self.hitUnits = ko.computed(function() {
                        if (self.units) {
                                return ko.utils.arrayFilter(self.units(), function(unit) {
                                        return unit().hit();
                                });
                        }
                });

                self.retreatableCav = ko.computed(function() {
                        return ko.utils.arrayFilter(self.units(), function(unit) {
                                return unit().newDamage() > 0 && unit().kind == "cavalry"
                        });
                });

                self.myUpdate = function(data) {
                        if (!data) {
                                return
                        }
                        if (self.hit() != data.hits) {
                                self.hit(data.hits);
                        }
                        if (data.commandBlock) {
                                if (self.commandBlock()) {
                                        self.commandBlock().myUpdate(data.commandBlock);
                                } else {
                                        self.commandBlock(new commandBlockModel(data.commandBlock))
                                }
                        } else {
                                self.commandBlock(null);
                        }

                        if (self.orderDiscs.changed(data.orderDiscs.data)) {
                                self.orderDiscs(data.orderDiscs.data);
                        }
                        self.units.myUpdate(data.units);
                };
        }

        function discModel(data) {
                var self = this;
                self.side = function(d) {
                        switch (d) {
                                case 1:
                                case 2:
                                        return "confederate";
                                case 3:
                                case 4:
                                        return "union";
                                default:
                                        return "none";
                        }
                }(data);

                self.forced = function(d) { return d == 4 }(data);

                self.double = function(d) { return d == 2 }(data);

                self.link = function(d) {
                        if (self.forced) {
                                return "/images/gettysburg/" + self.side + "-forced-pass-disc.png";
                        } else if (self.double) {
                                return "/images/gettysburg/" + self.side + "-double-order-disc.png";
                        } else {
                                return "/images/gettysburg/" + self.side + "-order-disc.png";
                        }
                }();
        }

        function initDiscs(data) {
                var discs = [];
                for (var i = 0; i < data.confederateSingle; i++) {
                        discs.push(new discModel(1));
                }
                for (var i = 0; i < data.confederateDouble; i++) {
                        discs.push(new discModel(2));
                }
                for (var i = 0; i < data.unionSingle; i++) {
                        discs.push(new discModel(3));
                }
                for (var i = 0; i < data.unionForcedPass; i++) {
                        discs.push(new discModel(4));
                }
                return discs;
        }

        function playerModel(data) {
                var self = this;

                self.id = data.p.id;
                self.gravatar = data.p.gravatar;
                self.link = data.p.link;
                self.log = ko.observableArray(data.log);

                self.toCommandBlocks = function(data) {
                        var blocks = [];
                        if (!data) {
                                return [];
                        }
                        for (var i = 0; i < data.length; i++) {
                                blocks.push(ko.observable(new commandBlockModel(data[i])));
                        }
                        return blocks;
                };
                self.commandBlocks = ko.observableArray(self.toCommandBlocks(data.commandBlocks));
                self.uCommandBlocks = ko.observableArray(self.toCommandBlocks(data.uCommandBlocks));

                self.orderDiscs = ko.observable(data.orderDiscs.data);
                self.uOrderDiscs = ko.observable(data.uOrderDiscs.data);
        }

        function commandBlockModel(data) {
                var self = this;

                self.side = data.side;
                self.value = data.value;
                self.faceUp = ko.observable(data.faceUp);

                self.link = ko.computed(function() {
                        if (self.faceUp() && self.value) {
                                return "/images/gettysburg/" + side(self.side) + "-command-block-" + self.value + ".png";
                        } else {
                                return "/images/gettysburg/" + side(self.side) + "-command-block-blank.png"
                        }
                });

                self.myUpdate = function(data) {
                        if (data && self.faceUp() != data.faceUp) {
                                self.faceUp(data.faceUp);
                        }
                }
        }


        return {
                create: function viewModel(data) { 
                        var self = this;

                        self.id = data.id;
                        self.title = ko.observable(data.title);
                        self.round = ko.observable(data.round);
                        self.phase = ko.observable(data.phase);
                        self.phaseName = ko.computed(function() {
                                return phaseNames[self.phase()]
                        });
                        self.subPhase = ko.observable(data.subPhase);
                        self.subPhaseName = ko.computed(function() {
                                return subPhaseNames[self.subPhase()]
                        });
                        self.status = ko.observable(data.status);
                        self.notices = ko.observableArray(data.notices);
                        self.errors = ko.observableArray(data.errors);
                        self.globalMap = ko.observable(data.globalMap);
                        self.globalMap.myUpdate = function(data) {
                                if (self.globalMap.changed(data)) {
                                        self.globalMap(data);
                                }
                        }
                        self.globalMap.changed = function(data) {
                                if (self.globalMap().zIndex != data.zIndex) {
                                        return true;
                                }
                                if (self.globalMap().enableOnClick != data.enableOnClick) {
                                        return true;
                                }
                                if (self.globalMap().selectedArea != data.selectedArea) {
                                        return true;
                                }
                                for (var i = 0; i < data.areaMap.length; i++) {
                                        if (self.globalMap().areaMap[i].isSelectable != data.areaMap[i].isSelectable) {
                                                return true;
                                        }
                                }
                                return false;
                        }

                        self.turn = ko.observable(data.turn);
                        self.turnName = data.turnName;
                        self.turnClass = function() {
                                if (self.turnName) {
                                        return self.turnName.replace(/ /g,'-').toLowerCase();
                                } else {
                                        return "none";
                                }
                        }();

                        self.toPlayers = function(data) {
                                var players = [];
                                if (data) {
                                        for (var i = 0; i < data.length; i++) {
                                                players.push(ko.observable(new playerModel(data[i])));
                                        }
                                }
                                return players;
                        };
                        self.players = ko.observableArray(self.toPlayers(data.players));
                        self.players.myUpdate = function(data) {
                                for (var i = 0; i < self.players.length; i++) {
                                        self.players[i](new playerModel(data[i]));
                                }
                        }

                        self.confederatePlayer = ko.computed(function() {
                                return self.players()[0];
                        });

                        self.unionPlayer = ko.computed(function() {
                                return self.players()[1];
                        });


                        self.period = data.period;
                        self.periodName = function() {
                                switch (self.period) {
                                        case 1:
                                                return "Morning";
                                        case 2:
                                                return "Mid-Day";
                                        case 3:
                                                return "Afternoon"; 
                                        case 4:
                                                return "Night";
                                        default:
                                                return "None";
                                }
                        }();

                        self.cpUserIndices = data.cpUserIndices ? ko.observable(data.cpUserIndices[0]) : ko.observable(-1);
                        self.currentPlayer = ko.computed(function() {
                                return ko.utils.arrayFirst(self.players(), function(player) {
                                        return player().id == self.cpUserIndices();
                                });
                        });

                        self.winnerIndices = data.winnerIndices ? ko.observable(data.winnerIndices[0]) : ko.observable(-1);

                        self.header = ko.computed(function() {
                                return {
                                        title: self.title(),
                                        turn: self.turn(),
                                        round: self.round(),
                                        phase: self.phase(),
                                        subPhase: self.subPhase(),
                                        cpUserIndices: self.cpUserIndices(),
                                        winnerIndices: self.winnerIndices(),
                                }
                        });
                        self.label = ko.observable('dialog test');
                        self.isOpen = ko.observable(false);
                        self.open = function() {
                                self.isOpen(true);   
                        };
                        self.close = function() {
                                self.isOpen(false);   
                        };

                        self.toAreas = function(data) {
                                var areas = [];
                                if (data) {
                                        for (var i = 0; i < data.length; i++) {
                                                areas.push(new areaModel(data[i]));
                                        }
                                }
                                return areas;
                        };
                        self.areas = self.toAreas(data.areas);
                        self.areas.myUpdate = function(data) {
                                for (var i = 0; i < self.areas.length; i++) {
                                        self.areas[i].myUpdate(data[i]);
                                }
                        };

                        self.selectedAreaID = ko.observable(data.selectedAreaID);
                        self.selectedArea = ko.computed(function() {
                                return self.areas[self.selectedAreaID()];
                        });
                        self.confederateEliteInfantryCubeCount = ko.observable(data.cubes.data.confederateEliteInfantry);
                        self.confederateRegularInfantryCubeCount = ko.observable(data.cubes.data.confederateRegularInfantry);
                        self.unionInferiorInfantryCubeCount = ko.observable(data.cubes.data.unionInferiorInfantry);
                        self.unionRegularInfantryCubeCount = ko.observable(data.cubes.data.unionRegularInfantry);
                        self.unionEliteInfantryCubeCount = ko.observable(data.cubes.data.unionEliteInfantry);
                        self.unionRegularCavalryCubeCount = ko.observable(data.cubes.data.unionRegularCavalry);

                        self.myUpdate = function(data) {
                                if (!data) {
                                        return
                                }
                                if (self.title() != data.title) {
                                        self.title(data.title)
                                }
                                if (self.turn() != data.turn) {
                                        self.turn(data.turn)
                                }
                                if (self.phase() != data.phase) {
                                        self.phase(data.phase)
                                }
                                if (self.subPhase() != data.subPhase) {
                                        self.subPhase(data.subPhase)
                                }
                                if (self.round() != data.round) {
                                        self.round(data.round)
                                }
                                if (self.status() != data.status) {
                                        self.status(data.status)
                                }
                                if (self.selectedAreaID() != data.selectedAreaID) {
                                        self.selectedAreaID(data.selectedAreaID)
                                }

                                self.notices(data.notices);
                                self.errors(data.errors);
                                self.globalMap.myUpdate(data.globalMap);
                                self.areas.myUpdate(data.areas);
                                self.players.myUpdate(data.players);
                        }
                },
                put: function(json) {
                        var $ajaxTarget = window.location.href;
                        $.ajax({
                                url: $ajaxTarget,
                                data: json, 
                                dataType: "json",
                                success: function(data) {
                                        gettysburg.myUpdate(data);
                                },
                                type: "PUT",
                                error: function(data) {
                                        alert("Error");
                                        window.location.href = $ajaxTarget;
                                }
                        }); 
                },
                clickHandler: function(data) {
                        var $ajaxTarget = window.location.href;
                        if (data.selected) {
                                $.ajax({
                                        url: $ajaxTarget,
                                        data: { "area" : data.key, "action" : "select-area"}, 
                                        dataType: "json",
                                        success: function(data) {
                                                gettysburg.myUpdate(data);
                                        },
                                        type: "PUT",
                                        error: function(data) {
                                                window.location.href = target;
                                        }
                                }); 
                        }
                },
                myUpdate: function(data) {
                        var initData = data,
                        dialogData = initData.dialog;

                        delete initData.dialog;

                        viewModel.myUpdate(initData);

                        if ($("#dialog").length > 0) {
                                ko.unapplyBindings($("#dialog"), true);
                        }
                        $("#board").append("<div id='dialog' data-bind='jqDialog: $root, dialogVisible: true'>" + dialogData.content + "</div>");
                        delete dialogData.content;
                        dialogViewModel = new myDialogViewModel(dialogData);
                        ko.applyBindings(dialogViewModel, $("#dialog")[0]);
                },
                assignDamage: function(element, valueAccessor, allBindingsAccessor, vm, bindingContext) {
                        var value = valueAccessor(),
                        change = $(element).spinner("value") - value(),
                        unit = bindingContext.$data,
                        key = unit.side + unit.rank.ucfirst() + unit.kind.ucfirst() + "CubeCount";
                        var supply = viewModel[key];

                        value($(element).spinner("value"));
                        supply(supply() - change);
                },

        }
}();

var myDialogViewModel = function(data) {
        var self = this;

        self.initialized = ko.observable(data.initialized);
        self.title = ko.observable(data.title);
        self.open = ko.observable(data.open);
        self.minWidth = ko.observable(data.minWidth);
        self.selectedArea = ko.observable(viewModel.selectedArea());
}


var areas = $("#areas");
var units = $(".unit");

// A single click handler for all 'img.unit' elements
$("#areas").on("click", "img.unit", function() {
        var context = ko.contextFor(this);
        var index = context.$index;
        var areaID = "area-" + context.$parent.id;
        var $ajaxTarget = window.location.href;
        var unitID = context.$data.id;

        $.ajax({
                url: $ajaxTarget,
                data: { "area": areaID, "action": "select-unit", "unit": unitID }, 
                dataType: "json",
                success: function(data) {
                        gettysburg.Update(data);
                },
                type: "PUT",
                error: function(data) {
                        window.location.href = $ajaxTarget;
                },
        }); 
});

ko.bindingHandlers.imageMapster = {
        init: function(element) {
                $(element).mapster({
                        scaleMap: true,
                fill: false,
                stroke: true,
                strokeWidth: 10,
                render_highlight: {
                        strokeOpacity: 0.75
                },
                render_select: {
                        strokeOpacity: 1.0
                },
                strokeColor: 'ffff00',
                onClick: gettysburg.clickHandler,
                mapKey: 'data-key',
                singleSelect: true,
                isDeselectable: false,
                showToolTip: true,
                });
        },
        myUpdate: function(element, valueAccessor) {
                var global = ko.unwrap(valueAccessor());
                $(element).mapster('set_options', { areas : ko.toJS(global.areaMap) });
                $(element).mapster('set_options', { onClick: gettysburg.clickHandler });
                if (!global.enableOnClick) {
                        $(element).mapster('set_options', { onClick: null });
                } else {
                        $(element).mapster('set_options', { onClick: gettysburg.clickHandler });
                }
                $(element).zIndex(global.zIndex);
                if (global.selectedArea != "none") {
                        $(element).mapster('set', true, global.selectedArea);
                }
        },
}

ko.bindingHandlers.jqButton = {
        init: function(element) {
                $(element).button();
        }
};

ko.bindingHandlers.jqTabs = {
        init: function(element) {
                $(element).tabs({ selected: 0});
        }
};

ko.bindingHandlers.jqSpinner = {
        init: function(element, valueAccessor, allBindingsAccessor, viewModel, bindingContext) {
                var value = valueAccessor(),
                min = 0,
                max = 100,
                $el = $(element);
                if (allBindingsAccessor.has("min")) {
                        min = allBindingsAccessor.get("min");
                }
                if (allBindingsAccessor.has("max")) {
                        max = allBindingsAccessor.get("max");
                }
                $el.css({width: "20px"}).spinner({min: min, max: max});
                $el.spinner("value", value());
                $el.on( "spinstop", function( event, ui ) {
                        if (allBindingsAccessor.has("update")) {
                                var myUpdate = allBindingsAccessor.get("update");
                                myUpdate(element, valueAccessor, allBindingsAccessor, viewModel, bindingContext);
                        } else {
                                value($el.spinner("value"));
                        }
                });
        },
        update: function(element, valueAccessor, allBindingsAccessor, viewModel, bindingContext) {
                if (allBindingsAccessor.has("update")) {
                        var myUpdate = allBindingsAccessor.get("update");
                        myUpdate(element, valueAccessor, allBindingsAccessor, viewModel, bindingContext);
                } else {
                        var value = valueAccessor();
                        value($(element).spinner("value"));
                }
        }
}

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

function positionDialog(dialog) {
        var boardTop = $("#board").offset().top;
        var windowTop = $(window).scrollTop();
        var topOffset = 0;
        if (windowTop < boardTop) {
                topOffset = boardTop - windowTop;
        }
        dialog.dialog("option", "position", {
                my: "right top",
                at: "right top+" + topOffset,
                of: window,
                within: "#board", 
                collision: "fit",
        }
        );
}

$(function() {
        $(window).scroll(function () {
                var $dialog = $("#dialog");
                if ($dialog.length > 0 && !$dialog.dialog("option", "modal") && $dialog.dialog("isOpen")) {
                        positionDialog($dialog);
                };
        });
});

ko.bindingHandlers.stopBinding = {
        init: function() {
                return { controlsDescendantBindings: true };
        }
};

ko.virtualElements.allowedBindings.stopBinding = true;

// accepts jQuery node and remove boolean
ko.unapplyBindings = function ($node, remove) {
        // unbind events
        $node.find("*").each(function () {
                $(this).unbind();
        });

        // Remove KO subscriptions and references
        if (remove) {
                ko.removeNode($node[0]);
        } else {
                ko.cleanNode($node[0]);
        }
};

ko.bindingHandlers.dataTable = {
        init: function(element) {
                $(element).dataTable( {
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
        }
}

