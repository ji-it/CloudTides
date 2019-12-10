import React from "react";
import PropTypes from "prop-types";
import {
    Row,
    Col,
    FormSelect,
    Card,
    CardHeader,
    CardBody,
    CardFooter
} from "shards-react";

import Chart from "../../utils/chart";

class UsageResource extends React.Component {
    constructor(props) {
        super(props);
        this.state = {};
        this.canvasRef = React.createRef();
    }

    componentDidMount() {
        const chartConfig = {
            type: "doughnut",
            data: this.props.chartData,
            options: {
                ...{
                    // tooltips: {
                    //     custom: false,
                    //     mode: "index",
                    //     position: "nearest"
                    // },
                    cutoutPercentage: 70,
                    // circumference: Math.PI,
                    // rotation: Math.PI,
                },
                ...this.props.chartOptions
            }
        };
        this.chart = new Chart(this.canvasRef.current, chartConfig);
    }

    render() {
        const {title, o_data} = this.props;
        const {resource} = o_data;
        let hosts = 0, vms = 0, contributing = 0;
        if (resource) {
            hosts = resource.hosts;
            vms = resource.vms;
            contributing = resource.contributing;
            const nf = new Intl.NumberFormat('en-US', {
                style: 'percent',
            });
            const percentage = nf.format(contributing / vms);
            const divide = contributing / vms;
            this.props.chartData.datasets[0].data = [divide, 1 - divide];
            this.props.chartOptions.elements.center.text = percentage + ' used';
            const chartConfig = {
                data: this.props.chartData,
                options: {
                    ...{
                        cutoutPercentage: 70,
                    },
                    ...this.props.chartOptions
                }
            };
            this.chart && (this.chart.data = chartConfig.data) && (this.chart.options = chartConfig.options) && this.chart.update()
            // this.state.chart = new Chart(this.canvasRef.current, chartConfig);
        }
        return (
            <Card small className="h-100">
                <CardHeader className="border-bottom">
                    <h6 className="m-0">{title}</h6>
                </CardHeader>
                <CardBody className="d-flex py-0">
                    <canvas
                        height="160"
                        ref={this.canvasRef}
                        className="m-auto"
                    />
                </CardBody>
                <CardFooter className="">
                    <Row>
                        <Col>
                            <span className="d-block small"><b>{hosts}</b> hosts</span>
                            <span className="d-block small"><b>{vms}</b> VMs</span>
                        </Col>
                        <Col className="text-right small">
                            <div className="bottom-aligner"></div>
                            <span><b>{contributing}</b> contributing</span>
                        </Col>
                    </Row>
                </CardFooter>
            </Card>
        );
    }
}

UsageResource.propTypes = {
    /**
     * The component's title.
     */
    title: PropTypes.string,
    data: PropTypes.object,
    /**
     * The chart config object.
     */
    chartConfig: PropTypes.object,
    /**
     * The Chart.js options.
     */
    chartOptions: PropTypes.object,
    /**
     * The chart data.
     */
    chartData: PropTypes.object
};

UsageResource.defaultProps = {
    title: "Resource Usage",
    data: {},
    chartData: {
        datasets: [
            {
                data: [0, 0],
                backgroundColor: [
                    "rgba(0,61,255,1.0)",
                    "rgba(229,228,234,1.0)",
                ]
            }
        ],
    },
    chartOptions: {
        events: [],
        elements: {
            arc: {
                roundedCornersFor: 0
            },
            center: {
                // the longest text that could appear in the center
                maxText: '100% used',
                text: '67% used',
                fontColor: '#565656',
                fontFamily: "'Helvetica Neue', 'Helvetica', 'Arial', sans-serif",
                fontStyle: 'normal',
                minFontSize: 1,
                maxFontSize: 30,
            }
        }
    }
};

export default UsageResource;
