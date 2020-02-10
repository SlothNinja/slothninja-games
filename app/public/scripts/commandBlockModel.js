define(['knockout-3.2.0'], function(ko) {
        return function commandBlockModel(data) {
                var self = this;

                self.sideInt = ko.observable(data.sideInt);
                self.side = ko.computed(function() {
                        if (self.sideInt() === 2) {
                                return "confederate";
                        }
                        return "union";
                });
                self.value = data.value;
                self.faceUp = ko.observable(data.faceUp);

                self.link = ko.computed(function() {
                        if (self.faceUp() && self.value) {
                                return "/images/gettysburg/" + self.side() + "-command-block-" + self.value + ".png";
                        } else {
                                return "/images/gettysburg/" + self.side() + "-command-block-blank.png"
                        }
                });

                self.myUpdate = function(data) {
                        if (data && self.faceUp() != data.faceUp) {
                                self.faceUp(data.faceUp);
                        }
                }
        }
})
