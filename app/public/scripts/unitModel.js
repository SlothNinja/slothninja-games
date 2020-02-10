define(['knockout-3.2.0'], function(ko) {
        return function unitModel(data) {
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
                self.selectable = ko.observable(data.selectable);
                self.selected = ko.observable(data.selected);
                self.fireFromID = ko.observable(data.fireFromID);
                self.targetAreaID = ko.observable(data.targetAreaID);
                self.targetUnits = ko.observable(data.targetUnits);

                self.link = function() {
                        return "/images/gettysburg/" + self.side + "-" + self.rank + "-" + self.kind + ".png";
                };

                self.totalDamage = ko.computed(function() {
                        return self.damage() + self.newDamage();
                });

                //                self.cubes = ko.computed(function() {
                //                        var cubes = [];
                //                        for (var i = 0; i < self.totalDamage(); i++) {
                //                                cubes.push(new cubeModel({side: self.side, rank: self.rank, kind: self.kind}));
                //                        }
                //                        return cubes
                //                });
                //

                self.cubelink = function() {
                        if (self.kind == 'artillery') {
                                return ''
                        } else {
                                return "/images/gettysburg/" + self.side + "-" + self.rank + "-" + self.kind + "-cube.png";
                        }
                }();

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
                        if (self.selectable() != data.selectable) {
                                self.selectable(data.selectable)
                        }
                        if (self.selected() != data.selected) {
                                self.selected(data.selected)
                        }
                        if (self.fireFromID() != data.fireFromID) {
                                self.fireFromID(data.fireFromID)
                        }
                        if (self.targetAreaID() != data.targetAreaID) {
                                self.targetAreaID(data.targetAreaID)
                        }
                        if (self.targetUnits() != data.targetUnits) {
                                self.targetUnits(data.targetUnits)
                        }
                }
        }
})
