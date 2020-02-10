define(['knockout-3.2.0', 'commandBlockModel'], function(ko, commandBlockModel) {
        return function playerModel(data) {
                var self = this;

                self.id = data.p.id;
                self.side = data.side;
                self.gravatar = data.p.gravatar;
                self.link = data.p.link;
                self.log = ko.observable(data.log);
                self.performedAction = ko.observable(data.p.performedAction);

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

                self.myUpdate = function(data) {
                        if (!data) {
                                return
                        };

                        if (data.commandBlocks.length != self.commandBlocks().length) {
                                self.commandBlocks(self.toCommandBlocks(data.commandBlocks));
                        };

                        if (data.uCommandBlocks.length != self.uCommandBlocks().length) {
                                self.uCommandBlocks(self.toCommandBlocks(data.uCommandBlocks));
                        };

                        if (data.orderDiscs.data != self.orderDiscs()) {
                                self.orderDiscs(data.orderDiscs.data)
                        }

                        if (data.uOrderDiscs.data != self.uOrderDiscs()) {
                                self.uOrderDiscs(data.uOrderDiscs.data)
                        }

                        if (data.log != self.log()) {
                                self.log(data.log)
                        }
                                                
//                        if (self.newDamage() != data.newDamage) {
//                                self.newDamage(data.newDamage);
//                        }
//                        if (data.commandBlock) {
//                                if (self.commandBlock()) {
//                                        self.commandBlock().myUpdate(data.commandBlock);
//                                } else {
//                                        self.commandBlock(new commandBlockModel(data.commandBlock))
//                                }
//                        } else {
//                                self.commandBlock(null);
//                        }
//
//                        if (self.selectable() != data.selectable) {
//                                self.selectable(data.selectable);
//                        };
//
//                        if (self.artilleryFireMarker() != data.artilleryFireMarker) {
//                                self.artilleryFireMarker(data.artilleryFireMarker)
//                        };
//
//                        if (self.disrupted() != data.disrupted) {
//                                self.disrupted(data.disrupted)
//                        };
//
//                        self.orderDiscs(data.orderDiscs);
//                        self.units.myUpdate(data.units);
                };
        }
})
