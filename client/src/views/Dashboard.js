import React from "react";
import PropTypes from "prop-types";
import {Container, Row, Col} from "shards-react";

import PageTitle from "./../components/common/PageTitle";
import Resources from "./../components/dashboard/Resources";
import UsageResource from "../components/dashboard/UsageResource";
import ContributionCard from "../components/dashboard/ContributionCard";
import PowerCard from "../components/dashboard/PowerCard";
import WorkloadsCard from "../components/dashboard/WorkloadsCard";
import Store from "../flux/store";
import {Actions} from "../flux";


class Dashboard extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            tableData: Store.getResourceTableData(),
            policies: Store.getPoliciesTableData(),
            overview: Store.getOverviewData(),
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
        Actions.getResources();
        Actions.getPolicies();
        Actions.getOverview();
        this.timer1 = Actions.getResources(true);
        this.timer2 = Actions.getOverview(true)
    }

    onChange() {
        this.setState({
            ...this.state,
            tableData: Store.getResourceTableData(),
            policies: Store.getPoliciesTableData(),
            overview: Store.getOverviewData(),
        });
    }

    render() {
        const {tableData: resourceData, policies: policiesData, overview: overviewData} = this.state;
        return (
            <Container fluid className="main-content-container px-4">
                {/* Page Header */}
                <Row noGutters className="page-header py-4">
                    <PageTitle title="Dashboard" subtitle="Overview" className="text-sm-left mb-3"/>
                </Row>
                <Row>
                    {/* Resources CPU Usage */}
                    <Col lg="3" md="4" sm="12" className="mb-4">
                        <UsageResource o_data={overviewData}/>
                    </Col>

                    {/* Contribution */}
                    <Col lg="2" md="4" sm="12" className="mb-4">
                        <ContributionCard data={overviewData}/>
                    </Col>

                    {/* Power */}
                    <Col lg="2" md="4" sm="12" className="mb-4">
                        <PowerCard data={overviewData}/>
                    </Col>

                    {/* Workloads */}
                    <Col lg="5" md="12" sm="12" className="mb-4">
                        <WorkloadsCard data={overviewData}/>
                    </Col>
                </Row>
                <Row>
                    {/* Resources */}
                    <Col className="col-lg mb-4">
                        <Resources data={resourceData} policiesData={policiesData}/>
                    </Col>
                </Row>
            </Container>
        )
    }
}

Dashboard.propTypes = {};

Dashboard.defaultProps = {};

export default Dashboard;
