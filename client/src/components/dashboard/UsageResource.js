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

        new Chart(this.canvasRef.current, chartConfig);
    }

    render() {
        const {title} = this.props;
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
                            <span className="d-block small"><b>208</b> hosts</span>
                            <span className="d-block small"><b>429</b> VMs</span>
                        </Col>
                        <Col className="text-right small">
                            <div className="bottom-aligner"></div>
                            <span><b>20</b> idle</span>
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
    chartData: {
        datasets: [
            {
                data: [68.3, 31.7],
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
