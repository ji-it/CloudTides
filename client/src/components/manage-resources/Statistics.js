import React from "react";
import PropTypes from "prop-types";
import {Row, Col, Card, CardHeader, CardBody, Button} from "shards-react";

import RangeDatePicker from "../common/RangeDatePicker";
import Chart from "../../utils/chart";
import _ from 'lodash'

class Statistics extends React.Component {
    constructor(props) {
        super(props);

        this.canvasRef = React.createRef();
    }

    componentDidMount() {
        const chartOptions = {
            ...{
                responsive: true,
                legend: {
                    position: "top"
                },
                elements: {
                    line: {
                        // A higher value makes the line look skewed at this ratio.
                        tension: 0.3
                    },
                    point: {
                        radius: 0
                    }
                },
                scales: {
                    xAxes: [
                        {
                            gridLines: false,
                            ticks: {
                                callback(tick, index) {
                                    // Jump every 7 values on the X axis labels to avoid clutter.
                                    return index % 7 !== 0 ? "" : tick;
                                }
                            }
                        }
                    ],
                    yAxes: [
                        {
                            ticks: {
                                suggestedMax: 45,
                                callback(tick) {
                                    if (tick === 0) {
                                        return tick;
                                    }
                                    // Format the amounts using Ks for thousands.
                                    return tick > 999 ? `${(tick / 1000).toFixed(1)}K` : tick;
                                }
                            }
                        }
                    ]
                },
                hover: {
                    mode: "nearest",
                    intersect: false
                },
                tooltips: {
                    custom: false,
                    mode: "nearest",
                    intersect: false
                }
            },
            ...this.props.chartOptions
        };

        const StatsOverview = new Chart(this.canvasRef.current, {
            type: "LineWithLine",
            data: this.props.chartData,
            options: chartOptions
        });

        // They can still be triggered on hover.
        const buoMeta = StatsOverview.getDatasetMeta(0);
        buoMeta.data.length > 0 && (buoMeta.data[0]._model.radius = 0) &&
        (buoMeta.data[this.props.chartData.datasets[0].data.length - 1]._model.radius = 0);

        // Render the chart.
        StatsOverview.render();
        this.chart = StatsOverview
    }

    render() {
        const {title, data} = this.props;
        if (!_.isEmpty(data)) {
            const key = Object.keys(data)[0];
            const labels = data[key]["time"];
            const ram = data[key]["ram"];
            const cpu = data[key]["cpu"];
            this.chart && (this.chart.data.datasets[1].data = ram) && (this.chart.data.datasets[0].data =
                cpu) && (this.chart.data.labels = labels) && this.chart.update();
        }


        // this.chart && (this.chart.data = chartConfig.data) && (this.chart.options = chartConfig.options) && this.chart.update()
        return (
            <Card small className="h-100">
                <CardHeader className="m-2 mb-0">
                    <div>
                        <div style={{display: "inline-block"}}>
                            <h6 className="m-0 font-weight-bold">{title}</h6>
                        </div>
                    </div>
                </CardHeader>
                <CardBody className="pt-0">
                    {/*<Row className="border-bottom py-2 bg-light">*/}
                    {/*    /!*<Col sm="6" className="d-flex mb-2 mb-sm-0">*!/*/}
                    {/*    /!*    <RangeDatePicker/>*!/*/}
                    {/*    /!*</Col>*!/*/}
                    {/*    <Col>*/}
                    {/*        <Button*/}
                    {/*            size="sm"*/}
                    {/*            className="d-flex btn-white ml-auto mr-auto ml-sm-auto mr-sm-0 mt-3 mt-sm-0"*/}
                    {/*        >*/}
                    {/*            View Full Report &rarr;*/}
                    {/*        </Button>*/}
                    {/*    </Col>*/}
                    {/*</Row>*/}
                    <canvas
                        height="120"
                        ref={this.canvasRef}
                        style={{maxWidth: "100% !important"}}
                    />
                </CardBody>
            </Card>
        );
    }
}

Statistics.propTypes = {
    /**
     * The component's title.
     */
    title: PropTypes.string,
    /**
     * The chart dataset.
     */
    chartData: PropTypes.object,
    /**
     * The Chart.js options.
     */
    chartOptions: PropTypes.object
};

Statistics.defaultProps = {
    title: "Statistics",
    chartData: {
        labels: Array.from(new Array(30), (_, i) => (i === 0 ? 1 : i)),
        datasets: [
            {
                label: "CPU",
                fill: "start",
                data: [0],
                backgroundColor: "rgba(0,123,255,0.1)",
                borderColor: "rgba(0,123,255,1)",
                pointBackgroundColor: "#ffffff",
                pointHoverBackgroundColor: "rgb(0,123,255)",
                borderWidth: 1.5,
                pointRadius: 0,
                pointHoverRadius: 3
            },
            {
                label: "RAM",
                fill: "start",
                data: [0],
                backgroundColor: "rgba(255,65,105,0.1)",
                borderColor: "rgba(255,65,105,1)",
                pointBackgroundColor: "#ffffff",
                pointHoverBackgroundColor: "rgba(255,65,105,1)",
                borderDash: [3, 3],
                borderWidth: 1,
                pointRadius: 0,
                pointHoverRadius: 2,
                pointBorderColor: "rgba(255,65,105,1)"
            }
        ]
    }
};

export default Statistics;
