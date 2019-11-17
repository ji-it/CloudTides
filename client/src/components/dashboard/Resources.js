import React from "react";
import PropTypes from "prop-types";
import {Row, Col, Card, CardHeader, CardBody, Button} from "shards-react";


class Resources extends React.Component {
    render() {
        const {title} = this.props;
        return (
            <Card small className="h-100">
                <CardHeader className="border-bottom">
                    <h6 className="m-0">{title}</h6>
                </CardHeader>
                <CardBody className="p-0 pb-3">
                    <table className="table mb-0">
                        <thead className="bg-light">
                        <tr>
                            <th scope="col" className="border-0">
                                Name
                            </th>
                            <th scope="col" className="border-0">
                                Status
                            </th>
                            <th scope="col" className="border-0">
                                IP Address
                            </th>
                            <th scope="col" className="border-0">
                                CPU
                            </th>
                            <th scope="col" className="border-0">
                                RAM
                            </th>
                            <th scope="col" className="border-0">
                                Disk
                            </th>
                            <th scope="col" className="border-0">
                                Jobs Done
                            </th>
                            <th scope="col" className="border-0">
                                Project
                            </th>
                            <th scope="col" className="border-0">
                                Active
                            </th>
                        </tr>
                        </thead>
                        <tbody>
                        <tr>
                            <td>New York Datacenter</td>
                            <td>Idle</td>
                            <td>192.168.0.1</td>
                            <td>60%</td>
                            <td>6/10GB</td>
                            <td>10/25GB</td>
                            <td>20</td>
                            <td>SETI@home</td>
                            <td>Yes</td>
                        </tr>
                        <tr>
                            <td>Old York Datacenter</td>
                            <td>Busy</td>
                            <td>192.168.0.1</td>
                            <td>60%</td>
                            <td>6/10GB</td>
                            <td>10/25GB</td>
                            <td>20</td>
                            <td>SETI@home</td>
                            <td>Yes</td>
                        </tr>
                        <tr>
                            <td>LA Datacenter</td>
                            <td>Contributing</td>
                            <td>192.168.0.1</td>
                            <td>60%</td>
                            <td>6/10GB</td>
                            <td>10/25GB</td>
                            <td>20</td>
                            <td>SETI@home</td>
                            <td>Yes</td>
                        </tr>
                        </tbody>
                    </table>
                </CardBody>
            </Card>
        );
    }
}

Resources.propTypes = {
    /**
     * The component's title.
     */
    title: PropTypes.string,
    /**
     * The table dataset.
     */
    tableData: PropTypes.object,
};

Resources.defaultProps = {
    title: "Resources",
    tableData: {
        labels: Array.from(new Array(30), (_, i) => (i === 0 ? 1 : i)),
        datasets: [
            {
                label: "Current Month",
                fill: "start",
                data: [
                    500,
                    800,
                    320,
                    180,
                    240,
                    320,
                    230,
                    650,
                    590,
                    1200,
                    750,
                    940,
                    1420,
                    1200,
                    960,
                    1450,
                    1820,
                    2800,
                    2102,
                    1920,
                    3920,
                    3202,
                    3140,
                    2800,
                    3200,
                    3200,
                    3400,
                    2910,
                    3100,
                    4250
                ],
                backgroundColor: "rgba(0,123,255,0.1)",
                borderColor: "rgba(0,123,255,1)",
                pointBackgroundColor: "#ffffff",
                pointHoverBackgroundColor: "rgb(0,123,255)",
                borderWidth: 1.5,
                pointRadius: 0,
                pointHoverRadius: 3
            },
            {
                label: "Past Month",
                fill: "start",
                data: [
                    380,
                    430,
                    120,
                    230,
                    410,
                    740,
                    472,
                    219,
                    391,
                    229,
                    400,
                    203,
                    301,
                    380,
                    291,
                    620,
                    700,
                    300,
                    630,
                    402,
                    320,
                    380,
                    289,
                    410,
                    300,
                    530,
                    630,
                    720,
                    780,
                    1200
                ],
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

export default Resources;
