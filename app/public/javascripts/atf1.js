$(function() {
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
        showToolTip: true,
        areas:  [{
                        key: "sippar", 
                        toolTip: "You can trade for one extra resource."
                },
                {
                        key: "babylon",
                        toolTip: "Take two extra armies when you start an empire."
                },
                {
                        key: "nippur",
                        toolTip: "+1 VP for each area you control in Sumer when you score your empire."
                },
                {
                        key: "shuruppak",
                        toolTip: "Attacker must expend one extra army to destroy one of your cities."
                },
                {
                        key: "uruk",
                        toolTip: "Place one worker in Scribes box."
                },
                {
                        key: "ur",
                        toolTip: "Take one extra textile cube during phase one."
                },
                {
                        key: "eridu",
                        toolTip: "Place one Worker in the Tool Maker box."
                },
                {
                        key: "dilmun",
                        toolTip: "Remove all Workers from this box when a Decline occurs."
                }]
        });
});
