import React from "react";
import PropTypes from "prop-types";
import {Container, Row, Col} from "shards-react";

import PageTitle from "./../components/common/PageTitle";
import Resources from "./../components/dashboard/Resources";
import UsageResource from "../components/dashboard/UsageResource";
import ContributionCard from "../components/dashboard/ContributionCard";
import PowerCard from "../components/dashboard/PowerCard";
import WorkloadsCard from "../components/dashboard/WorkloadsCard";

const Dashboard = () => (
    <Container fluid className="main-content-container px-4">
        {/* Page Header */}
        <Row noGutters className="page-header py-4">
            <PageTitle title="Dashboard" subtitle="Overview" className="text-sm-left mb-3"/>
        </Row>
        <Row>
            {/* Resources CPU Usage */}
            <Col lg="3" md="4" sm="12" className="mb-4">
                <UsageResource/>
            </Col>

            {/* Contribution */}
            <Col lg="2" md="4" sm="12" className="mb-4">
                <ContributionCard/>
            </Col>

            {/* Power */}
            <Col lg="2" md="4" sm="12" className="mb-4">
                <PowerCard/>
            </Col>

            {/* Workloads */}
            <Col lg="5" md="12" sm="12" className="mb-4">
                <WorkloadsCard/>
            </Col>
        </Row>
        <Row>
            {/* Resources */}
            <Col className="col-lg mb-4">
                <Resources/>
            </Col>
        </Row>
    </Container>
);

Dashboard.propTypes = {};

Dashboard.defaultProps = {};

export default Dashboard;
