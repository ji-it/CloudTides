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
        const value = "Yes"
        return (
            <Card small className="h-100">
                <CardBody className="mt-3 d-flex py-0">
                    <Row className="w-100">
                        <Col lg="3" md="6" className="mb-4">
                            <Row className="w-100 m-auto">
                                <Col lg="6" className="">
                                    <div className="text-center mb-3">
                                        <h5 className="mb-0"> 0 </h5>
                                        Hosts
                                    </div>
                                    <div className="text-center">
                                        <h5 className="mb-0"> 2/3 </h5>
                                        Clusters
                                    </div>
                                </Col>
                                <Col lg="6" className="">
                                    <div className="text-center mb-3">
                                        <h5 className="mb-0"> Idle </h5>
                                        Status
                                    </div>
                                    <div className="text-center">
                                        <h5 className="mb-0"> 17 </h5>
                                        VMs
                                    </div>
                                </Col>
                            </Row>
                        </Col>
                        <Col lg="5" md="6" className="m-auto">
                            <Row>
                                <Col sm="3" className="font-weight-bold text-right" style={{fontSize: "12px"}}>CPU</Col>
                                <Col sm="6"> <Progress theme="primary" value={30}/></Col>
                                <Col sm="3" className="small text-left" style={{fontSize: "11px"}}>15.32/48GHz</Col>
                            </Row>
                            <Row>
                                <Col sm="3" className="font-weight-bold text-right"
                                     style={{fontSize: "12px"}}>Memory</Col>
                                <Col sm="6"> <Progress theme="primary" value={8.8}/></Col>
                                <Col sm="3" className="small text-left" style={{fontSize: "11px"}}>2.83/32GB</Col>
                            </Row>
                            <Row>
                                <Col sm="3" className="font-weight-bold text-right"
                                     style={{fontSize: "12px"}}>Storage</Col>
                                <Col sm="6"> <Progress theme="primary" value={57}/></Col>
                                <Col sm="3" className="small text-left" style={{fontSize: "11px"}}>1.03/1.81TB</Col>
                            </Row>
                        </Col>
                        <Col lg="4" md="12" className="m-auto text-right">
                            <div className="mb-0"><h4 className="mb-0 d-inline">New York Datacenter2</h4>
                                <a href="">
                                    <i className="fa fa-ellipsis-v ml-2 text-reagent-gray"
                                       style={{verticalAlign: "text-top"}}></i>
                                </a>
                            </div>
                            <div className="small mt-0">IP <span className="font-weight-bold">10.11.16.98</span></div>
                            <div className="small">
                                <MuiThemeProvider theme={this.getMuiTheme()}>
                                    <FormControlLabel
                                        label="Allow to Contribute"
                                        labelPlacement="start"
                                        value={value ? "Yes" : "No"}
                                        control={
                                            <Switch color="primary" size="small" checked={true}
                                                    value={value ? "Yes" : "No"}/>
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
};

export default SummaryCard;
