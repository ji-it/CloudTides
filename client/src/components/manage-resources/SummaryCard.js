import React from "react";
import PropTypes from "prop-types";
import {
    Row,
    Col,
    Card,
    CardHeader,
    CardBody,
    Progress,
    Button
} from "shards-react";
import FormControlLabel from "@material-ui/core/FormControlLabel";
import Switch from "@material-ui/core/Switch";
import {createMuiTheme, MuiThemeProvider} from "@material-ui/core";

class SummaryCard extends React.Component {

    getMuiTheme = () => createMuiTheme({
        overrides: this.props.tableStyle
    });

    render() {
        let {data, r_data} = this.props;
        if (r_data)
            data = r_data;
        const nf = new Intl.NumberFormat('en-US', {
            maximumFractionDigits: 1,
        });
        return (
            <Card small className="h-100">
                <CardBody className="mt-3 d-flex py-0">
                    <Row className="w-100">
                        <Col lg="3" md="6" className="mb-4">
                            <Row className="w-100 m-auto">
                                <Col lg="6" className="">
                                    <div className="text-center mb-3">
                                        <h5 className="mb-0"> {data.total_vms} </h5>
                                        VMs
                                    </div>
                                    <div className="text-center">
                                        <h5 className="mb-0"> {data.active_vms}</h5>
                                        Contrib. VMs
                                    </div>
                                </Col>
                                <Col lg="6" className="">
                                    <div className="text-center mb-3">
                                        <h5 className="mb-0"> {data.status} </h5>
                                        Status
                                    </div>
                                    <div className="text-center">
                                        <h5 className="mb-0"> {data.policy_name} </h5>
                                        Policy
                                    </div>
                                </Col>
                            </Row>
                        </Col>
                        <Col lg="5" md="6" className="m-auto">
                            <Row>
                                <Col sm="3" className="font-weight-bold text-right"
                                     style={{fontSize: "12px"}}>CPU</Col>
                                <Col sm="6"> <Progress theme="primary" value={data.current_cpu / data.total_cpu * 100}/></Col>
                                <Col sm="3" className="small text-left" style={{fontSize: "11px"}}>
                                    {nf.format(data.current_cpu)}/{nf.format(data.total_cpu)}GHz
                                </Col>
                            </Row>
                            <Row>
                                <Col sm="3" className="font-weight-bold text-right"
                                     style={{fontSize: "12px"}}>Memory</Col>
                                <Col sm="6"> <Progress theme="primary" value={data.current_ram / data.total_ram * 100}/></Col>
                                <Col sm="3" className="small text-left" style={{fontSize: "11px"}}>
                                    {nf.format(data.current_ram)}/{nf.format(data.total_ram)}GB
                                </Col>
                            </Row>
                            {/*<Row>*/}
                            {/*    <Col sm="3" className="font-weight-bold text-right"*/}
                            {/*         style={{fontSize: "12px"}}>Storage</Col>*/}
                            {/*    <Col sm="6"> <Progress theme="primary" value={57}/></Col>*/}
                            {/*    <Col sm="3" className="small text-left" style={{fontSize: "11px"}}>1.03/1.81TB</Col>*/}
                            {/*</Row>*/}
                        </Col>
                        <Col lg="4" md="12" className="m-auto text-right">
                            <div className="mb-0"><h4 className="mb-0 d-inline">{data.datacenter}</h4>
                                <a href="">
                                    <i className="fa fa-ellipsis-v ml-2 text-reagent-gray"
                                       style={{verticalAlign: "text-top"}}></i>
                                </a>
                            </div>
                            <div className="small mt-0">IP <span className="font-weight-bold">{data.host_name}</span>
                            </div>
                            <div className="small">
                                <MuiThemeProvider theme={this.getMuiTheme()}>
                                    <FormControlLabel
                                        label="Allow to Contribute"
                                        labelPlacement="start"
                                        value={data.is_active ? "Yes" : "No"}
                                        control={
                                            <Switch color="primary" size="small" checked={data.is_active}
                                                    value={data.is_active ? "Yes" : "No"}/>
                                        }
                                        // onChange={}
                                    />
                                </MuiThemeProvider>
                            </div>
                            <div className="mt-4 mb-3">
                                <Button theme="danger">Remove Resource</Button>
                            </div>
                        </Col>
                    </Row>
                </CardBody>
            </Card>
        );
    }
}

SummaryCard.propTypes = {
    /**
     * The component's title.
     */
    title: PropTypes.string,
    data: PropTypes.object,
    tableStyle: PropTypes.object,
};

SummaryCard.defaultProps = {
    title: "Summary",
    tableStyle: {
        MuiFormControlLabel: {
            label: {
                fontSize: "14px",
                fontWeight: "500"
            }
        },
    },
    data: {
        id: 28,
        date_added: "2019-12-08T07:59:43.942741Z",
        host_name: "",
        status: "",
        policy_name: "",
        platform_type: "",
        datacenter: "",
        total_cpu: 0,
        total_ram: 0,
        total_disk: null,
        current_ram: 0,
        current_cpu: 0,
        is_active: false,
        total_jobs: 0,
        ram_percent: 0,
        job_completed: 0,
        monitored: false,
        cpu_percent: 0,
        total_vms: 0,
        active_vms: 0,
        last_deployed: "2019-12-10T17:15:03Z"
    }
};

export default SummaryCard;
