const ws = new WebSocket("ws://localhost:5055/ws")

var dataSize = 3000
var dataAverage = 1
var receivedData = []

ws.addEventListener('message', (event) => {

    var jsonObj = JSON.parse(event.data)

    receivedData.push(jsonObj)
    
    if (receivedData.length > dataSize) {
        receivedData = receivedData.slice(Math.floor(dataSize / 10), dataSize)
    }
})

function applyDataParameters() {
    var elSize = document.getElementById('data-size');
    var elAverage = document.getElementById('data-average');

    if (elSize != undefined && elAverage != undefined) {

        try {
            var size = new Number(elSize.value)
            var average = new Number(elAverage.value)
            
            if(size <= 0 || average <= 0 || size < average) {
                window.alert("Invalid input!")
            } else {
                dataSize = size
                dataAverage = average
            }

        } catch (error) {
            window.alert("Invalid input!")
            console.log(error)
        }

    }
}

google.charts.load('current', {'packages':['corechart']});
google.charts.setOnLoadCallback(generateDeviationGraph);

function generateDeviationGraph() {

    setInterval(() => {

        var graphData = []
        var averageDeviation = 0

        receivedData.forEach((v, i) => {

            if (i % dataAverage == 0) {

                if(averageDeviation == 0) {
                    graphData.push([i, v.deviation])
                } else {
                    graphData.push([i, averageDeviation])
                }

                averageDeviation = 0
            } else {

                if(averageDeviation == 0) {
                    averageDeviation = v.deviation
                } else {
                    averageDeviation = (averageDeviation + v.deviation) / 2
                }

            }

        })

        var data = google.visualization.arrayToDataTable([
            ['Index', 'Deviation'],
            ...graphData
          ]);
      
          var options = {
            title: 'Company Performance',
            curveType: 'function',
            legend: { position: 'bottom' }
          };
      
          var chart = new google.visualization.LineChart(document.getElementById('curve_chart'));
          chart.draw(data, options);
    }, 3000)
}

