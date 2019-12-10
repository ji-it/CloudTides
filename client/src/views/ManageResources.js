import React from "react";
import PropTypes from "prop-types";
import {Container, Row, Col} from "shards-react";

import PageTitle from "./../components/common/PageTitle";
import Statistics from "../components/manage-resources/Statistics";
import SummaryCard from "../components/manage-resources/SummaryCard";
import Store from "../flux/store";
import {Actions} from "../flux";
import VMTable from "../components/manage-resources/VMTable";


class ManageResources extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            resources: Store.getDetailedResourceTableData(),
            resourceIndex: Store.getResourceIndex(),
            vms: Store.getVMTableData(),
        };

        this.onChange = this.onChange.bind(this);
    }

    componentWillMount() {
        Store.addChangeListener(this.onChange);
    }

    componentWillUnmount() {
        Store.removeChangeListener(this.onChange);
        clearInterval(this.timer1);
        this.timer1 = null;
        clearInterval(this.timer2);
        this.timer2 = null;
    }

    componentDidMount() {
        Actions.getDetailedResources();
        Actions.getVMS();
        this.timer1 = Actions.getDetailedResources(true);
        this.timer2 = Actions.getVMS(true);
    }

    onChange() {
        this.setState({
            ...this.state,
            resources: Store.getDetailedResourceTableData(),
            resourceIndex: Store.getResourceIndex(),
            vms: Store.getVMTableData(),
        });
    }

    render() {
        let resourceInDisplay = null;
        const resIndex = this.state.resourceIndex;
        let {vms: vmsData} = this.state;
        if (this.state.resources.length > 0) {
            resourceInDisplay = this.state.resources[0];
        }
        if (vmsData.length > 0) {
            vmsData = vmsData[0];
        }
        if (resIndex && this.state.resources.length > 0) {
            resourceInDisplay = this.state.resources[resIndex];
            vmsData = vmsData[resIndex]
        }

        return (
            <Container fluid className="main-content-container px-4">
                {/* Page Header */}
                <Row noGutters className="page-header py-4">
                    <PageTitle title="" subtitle="Manage Resources" className="text-sm-left mb-3"/>
                </Row>
                <Row>
                    {/* Informaiton */}
                    <Col lg="12" className="mb-4">
                        <SummaryCard r_data={resourceInDisplay}/>
                    </Col>
                </Row>
                <Row>
                    {/* VMTable */}
                    <Col className="col-lg mb-4">
                        <VMTable data={vmsData}/>
                    </Col>
                </Row>
                <Row>
                    {/* Statistics */}
                    <Col className="col-lg mb-4">
                        <Statistics/>
                    </Col>
                </Row>
            </Container>
        )
    }
}

ManageResources.propTypes = {};

ManageResources.defaultProps = {};

export default ManageResources;
