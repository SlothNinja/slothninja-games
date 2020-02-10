$(function() {
        InitDeedTips();
        InitGoodsTips();
        InitShipTips();
        InitCityTips();
        InitClickMap();
        InitZoomButtons();
});

MyAjax.Update = function(data) {
        this.update(data);
        InitDeedTips();
        InitGoodsTips();
        InitShipTips();
        InitCityTips();
        InitClickMap();
        InitZoomButtons();
        $('table.strippedDialogDataTable').dataTable( {
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

function InitDeedTips() {
        var $tips = $("div[data-deedtip]");
        $tips.tooltip({items: "div[data-deedtip]", content: function(){
                var $province = $(this).data("deedtip").province;
                var $goods = $(this).data("deedtip").goods;
                var $capacity = $(this).data("deedtip").capacity;
                if ($capacity != "0 0 0") {
                        return "<div>"+$province+"</div><div>Type: "+$goods+"</div><div>Capacity: "+$capacity+"</div>";
                } else {
                        return "<div>"+$province+"</div><div>Type: "+$goods+"</div>";
                }
        }});
};

function InitGoodsTips() {
        var $tips = $("img[data-goodstip]");
        $tips.tooltip({items: "img[data-goodstip]", content: function(){
                var $owner = $(this).data("goodstip").owner;
                var $slot = $(this).data("goodstip").slot;
                return "<div>Owner: "+$owner+"</div><div>Slot: "+$slot+"</div>";
        }});
};

function InitShipTips() {
        var $tips = $("img[data-shiptip]");
        $tips.tooltip({items: "img[data-shiptip]", content: function(){
                var $owner = $(this).data("shiptip").owner;
                var $slot = $(this).data("shiptip").slot;
                var $hull = $(this).data("shiptip").hull;
                var $delivered = $(this).data("shiptip").delivered;
                return "<div>Owner: "+$owner+"</div><div>Slot: "+$slot+"</div><div>Hull: "+$hull+
                "</div><div>Delivered: "+$delivered+"</div>";
        }});
};

function InitCityTips() {
        var $tips = $("img[data-citytip]");
        $tips.tooltip({items: "img[data-citytip]", content: function(){
                var $tipString = "";
                var $province =$(this).data("citytip").province;
                if ( typeof $province !== 'undefined') {
                        $tipString += "<div><strong>"+$province+"</strong></div>";
                }
                var $rice =$(this).data("citytip").rice;
                if ( typeof $rice !== 'undefined') {
                        $tipString += "<div>Rice: "+$rice+"</div>";
                }
                var $spice = $(this).data("citytip").spice;
                if ( typeof $spice !== 'undefined') {
                        $tipString += "<div>Spice: "+$spice+"</div>";
                }
                var $rubber = $(this).data("citytip").rubber;
                if ( typeof $rubber !== 'undefined') {
                        $tipString += "<div>Rubber: "+$rubber+"</div>";
                }
                var $oil = $(this).data("citytip").oil;
                if ( typeof $oil !== 'undefined') {
                        $tipString += "<div>Oil: "+$oil+"</div>";
                }
                var $siapfaji = $(this).data("citytip").siapFaji;
                if ( typeof $siapfaji !== 'undefined') {
                        $tipString += "<div>Siap Faji: "+$siapfaji+"</div>";
                }
                return $tipString;
        }});
};

function InitZoomButtons() {
        var $current = $("#board").data("zoom");
        $("#zoom-all").off("click").on("click", function() {
                $("#map-container").css({"height": "561"}); 
                $("#board").removeClass("zoom-left zoom-center zoom-right");
                $("#board").data("zoom", "zoom-all");
                $("#board-image").mapster('resize', 1190);
                $("#map-content-box").css({"position": "absolute", "left": "0px"});
        });
        $("#zoom-left").off("click").on("click", function() {
                $("#map-container").css({"height": "1122"}); 
                $("#board").removeClass("zoom-center zoom-right").addClass("zoom-left");
                $("#board").data("zoom", "zoom-left");
                $("#board-image").mapster('resize', 2380);
                $("#map-content-box").css({"position": "absolute", "left": "0px"});
        });
        $("#zoom-center").off("click").on("click", function() {
                $("#map-container").css({"height": "1122"}); 
                $("#board").removeClass("zoom-left zoom-right").addClass("zoom-center");
                $("#board").data("zoom", "zoom-center");
                $("#board-image").mapster('resize', 2380);
                $("#map-content-box").css({"position": "absolute", "left": "-590px"});
        });
        $("#zoom-right").off("click").on("click", function() {
                $("#map-container").css({"height": "1122"}); 
                $("#board").removeClass("zoom-left zoom-center").addClass("zoom-right");
                $("#board").data("zoom", "zoom-right");
                $("#board-image").mapster('resize', 2380); 
                $("#map-content-box").css({"position": "absolute", "left": "-1180px"});
        });
        $("#"+$current).trigger("click");
};

//function InitZoomButtons() {
//        $("#zoom-all").on("click", function() {
//                UnbindClickMap();
//                $("#map-container").css({"height": "561"}); 
//                $("#board").removeClass("zoom-left zoom-right");
//                $("#map-content-box").animate({"position": "absolute", "left": "0px"});
//                $("#board-image").animate({"width": "1190px", "height": "561px"}, function() {
//                        InitClickMap();
//                });
//        });
//        $("#zoom-left").on("click", function() {
//                UnbindClickMap();
//                $("#map-container").css({"height": "1122"}); 
//                $("#board").removeClass("zoom-right").addClass("zoom-left");
//                $("#map-content-box").animate({"position": "absolute", "left": "0px"});
//                $("#board-image").animate({"width": "2380px", "height": "1122px"}, function() {
//                        InitClickMap();
//                });
//        });
//        $("#zoom-right").on("click", function() {
//                UnbindClickMap();
//                $("#map-container").css({"height": "1122"}); 
//                $("#board").removeClass("zoom-left").addClass("zoom-right");
//                $("#map-content-box").animate({"position": "absolute", "left": "-1180px"});
//                $("#board-image").animate({"width": "2380px", "height": "1122px"}, function() {
//                        InitClickMap();
//                });
//        });
//};

//function UnbindClickMap() {
//        $(".clickmap").mapster('unbind');
//};

function InitClickMap() {
        $(".clickmap").mapster({
                scaleMap: true,
        fill: false,
        stroke: true,
        strokeWidth: 3,
        strokeOpacity: 0.75,
        strokeColor: 'ffff00',
        onClick: mapClickHandler,
        mapKey: 'data-key',
        singleSelect: true,
        isDeselectable: false,
        areas:  [
        { key: "area-0",
                toolTip: "Area: 0",
        },
        { key: "area-1",
                toolTip: "Area: 1",
        },
        { key: "area-2",
                toolTip: "Area: 2",
        },
        { key: "area-3",
                toolTip: "Area: 3",
        },
        { key: "area-4",
                toolTip: "Area: 4",
        },
        { key: "area-5",
                toolTip: "Area: 5",
        },
        { key: "area-6",
                toolTip: "Area: 6",
        },
        { key: "area-7",
                toolTip: "Area: 7",
        },
        { key: "area-8",
                toolTip: "Area: 8",
        },
        { key: "area-9",
                toolTip: "Area: 9",
        },
        { key: "area-10",
                toolTip: "Area: 10",
        },
        { key: "area-11",
                toolTip: "Area: 11",
        },
        { key: "area-12",
                toolTip: "Area: 12",
        },
        { key: "area-13",
                toolTip: "Area: 13",
        },
        { key: "area-14",
                toolTip: "Area: 14",
        },
        { key: "area-15",
                toolTip: "Area: 15",
        },
        { key: "area-16",
                toolTip: "Area: 16",
        },
        { key: "area-17",
                toolTip: "Area: 17",
        },
        { key: "area-18",
                toolTip: "Area: 18",
        },
        { key: "area-19",
                toolTip: "Area: 19",
        },
        { key: "area-20",
                toolTip: "Area: 20",
        },
        { key: "area-21",
                toolTip: "Area: 21",
        },
        { key: "area-22",
                toolTip: "Area: 22",
        },
        { key: "area-23",
                toolTip: "Area: 23",
        },
        { key: "area-24",
                toolTip: "Area: 24",
        },
        { key: "area-25",
                toolTip: "Area: 25",
        },
        { key: "area-26",
                toolTip: "Area: 26",
        },
        { key: "area-27",
                toolTip: "Area: 27",
        },
        { key: "area-28",
                toolTip: "Area: 28",
        },
        { key: "area-29",
                toolTip: "Area: 29",
        },
        { key: "area-30",
                toolTip: "Area: 30",
        },
        { key: "area-31",
                toolTip: "Area: 31",
        },
        { key: "area-32",
                toolTip: "Area: 32",
        },
        { key: "area-33",
                toolTip: "Area: 33",
        },
        { key: "area-34",
                toolTip: "Area: 34",
        },
        { key: "area-35",
                toolTip: "Area: 35",
        },
        { key: "area-36",
                toolTip: "Area: 36",
        },
        { key: "area-37",
                toolTip: "Area: 37",
        },
        { key: "area-38",
                toolTip: "Area: 38",
        },
        { key: "area-39",
                toolTip: "Area: 39",
        },
        { key: "area-40",
                toolTip: "Area: 40",
        },
        { key: "area-41",
                toolTip: "Area: 41",
        },
        { key: "area-42",
                toolTip: "Area: 42",
        },
        { key: "area-43",
                toolTip: "Area: 43",
        },
        { key: "area-44",
                toolTip: "Area: 44",
        },
        { key: "area-45",
                toolTip: "Area: 45",
        },
        { key: "area-46",
                toolTip: "Area: 46",
        },
        { key: "area-47",
                toolTip: "Area: 47",
        },
        { key: "area-48",
                toolTip: "Area: 48",
        },
        { key: "area-49",
                toolTip: "Area: 49",
        },
        { key: "area-50",
                toolTip: "Area: 50",
        },
        { key: "area-51",
                toolTip: "Area: 51",
        },
        { key: "area-52",
                toolTip: "Area: 52",
        },
        { key: "area-53",
                toolTip: "Area: 53",
        },
        { key: "area-54",
                toolTip: "Area: 54",
        },
        { key: "area-55",
                toolTip: "Area: 55",
        },
        { key: "area-56",
                toolTip: "Area: 56",
        },
        { key: "area-57",
                toolTip: "Area: 57",
        },
        { key: "area-58",
                toolTip: "Area: 58",
        },
        { key: "area-59",
                toolTip: "Area: 59",
        },
        { key: "area-60",
                toolTip: "Area: 60",
        },
        { key: "area-61",
                toolTip: "Area: 61",
        },
        { key: "area-62",
                toolTip: "Area: 62",
        },
        { key: "area-63",
                toolTip: "Area: 63",
        },
        { key: "area-64",
                toolTip: "Area: 64",
        },
        { key: "area-65",
                toolTip: "Area: 65",
        },
        { key: "area-66",
                toolTip: "Area: 66",
        },
        { key: "area-67",
                toolTip: "Area: 67",
        },
        { key: "area-68",
                toolTip: "Area: 68",
        },
        { key: "area-69",
                toolTip: "Area: 69",
        },
        { key: "area-70",
                toolTip: "Area: 70",
        },
        { key: "area-71",
                toolTip: "Area: 71",
        },
        { key: "area-72",
                toolTip: "Area: 72",
        },
        { key: "area-73",
                toolTip: "Area: 73",
        },
        { key: "area-74",
                toolTip: "Area: 74",
        },
        { key: "area-75",
                toolTip: "Area: 75",
        },
        { key: "area-76",
                toolTip: "Area: 76",
        },
        { key: "area-77",
                toolTip: "Area: 77",
        },
        { key: "area-78",
                toolTip: "Area: 78",
        },
        { key: "area-79",
                toolTip: "Area: 79",
        },
        { key: "area-80",
                toolTip: "Area: 80",
        },
        { key: "area-81",
                toolTip: "Area: 81",
        },
        { key: "area-82",
                toolTip: "Area: 82",
        },
        { key: "area-83",
                toolTip: "Area: 83",
        },
        { key: "area-84",
                toolTip: "Area: 84",
        },
        { key: "area-85",
                toolTip: "Area: 85",
        },
        { key: "area-86",
                toolTip: "Area: 86",
        },
        { key: "area-87",
                toolTip: "Area: 87",
        },
        { key: "area-88",
                toolTip: "Area: 88",
        },
        { key: "area-89",
                toolTip: "Area: 89",
        },
        { key: "area-90",
                toolTip: "Area: 90",
        },
        { key: "area-91",
                toolTip: "Area: 91",
        },
        { key: "area-92",
                toolTip: "Area: 92",
        },
        { key: "area-93",
                toolTip: "Area: 93",
        },
        { key: "area-94",
                toolTip: "Area: 94",
        },
        { key: "area-95",
                toolTip: "Area: 95",
        },
        { key: "area-96",
                toolTip: "Area: 96",
        },
        { key: "area-97",
                toolTip: "Area: 97",
        },
        { key: "area-98",
                toolTip: "Area: 98",
        },
        { key: "area-99",
                toolTip: "Area: 99",
        },
        { key: "area-100",
                toolTip: "Area: 100",
        },
        { key: "area-101",
                toolTip: "Area: 101",
        },
        { key: "area-102",
                toolTip: "Area: 102",
        },
        { key: "area-103",
                toolTip: "Area: 103",
        },
        { key: "area-104",
                toolTip: "Area: 104",
        },
        { key: "area-105",
                toolTip: "Area: 105",
        },
        { key: "area-106",
                toolTip: "Area: 106",
        },
        { key: "area-107",
                toolTip: "Area: 107",
        },
        { key: "area-108",
                toolTip: "Area: 108",
        },
        { key: "area-109",
                toolTip: "Area: 109",
        },
        { key: "area-110",
                toolTip: "Area: 110",
        },
        { key: "area-111",
                toolTip: "Area: 111",
        },
        { key: "area-112",
                toolTip: "Area: 112",
        },
        { key: "area-113",
                toolTip: "Area: 113",
        },
        { key: "area-114",
                toolTip: "Area: 114",
        },
        { key: "area-115",
                toolTip: "Area: 115",
                strokeColor: '69FAFD',
        },
        { key: "area-116",
                toolTip: "Area: 116",
                strokeColor: '69FAFD',
        },
        { key: "area-117",
                toolTip: "Area: 117",
                strokeColor: '69FAFD',
        },
        { key: "area-118",
                toolTip: "Area: 118",
                strokeColor: '69FAFD',
        },
        { key: "area-119",
                toolTip: "Area: 119",
                strokeColor: '69FAFD',
        },
        { key: "area-120",
                toolTip: "Area: 120",
                strokeColor: '69FAFD',
        },
        { key: "area-121",
                toolTip: "Area: 121",
                strokeColor: '69FAFD',
        },
        { key: "area-122",
                toolTip: "Area: 122",
                strokeColor: '69FAFD',
        },
        { key: "area-123",
                toolTip: "Area: 123",
                strokeColor: '69FAFD',
        },
        { key: "area-124",
                toolTip: "Area: 124",
                strokeColor: '69FAFD',
        },
        { key: "area-125",
                toolTip: "Area: 125",
                strokeColor: '69FAFD',
        },
        { key: "area-126",
                toolTip: "Area: 126",
                strokeColor: '69FAFD',
        },
        { key: "area-127",
                toolTip: "Area: 127",
                strokeColor: '69FAFD',
        },
        { key: "area-128",
                toolTip: "Area: 128",
                strokeColor: '69FAFD',
        },
        { key: "area-129",
                toolTip: "Area: 129",
                strokeColor: '69FAFD',
        },
        { key: "area-130",
                toolTip: "Area: 130",
                strokeColor: '69FAFD',
        },
        { key: "area-131",
                toolTip: "Area: 131",
                strokeColor: '69FAFD',
        },
        { key: "area-132",
                toolTip: "Area: 132",
                strokeColor: '69FAFD',
        },
        { key: "area-133",
                toolTip: "Area: 133",
                strokeColor: '69FAFD',
        },
        { key: "area-134",
                toolTip: "Area: 134",
                strokeColor: '69FAFD',
        },
        { key: "area-135",
                toolTip: "Area: 135",
                strokeColor: '69FAFD',
        },
        ]
        });
};
