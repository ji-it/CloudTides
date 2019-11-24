import React from "react";
import PropTypes from "prop-types";
import {Container, Row, Col} from "shards-react";

import PageTitle from "./../components/common/PageTitle";
import Statistics from "../components/manage-resources/Statistics";
import SummaryCard from "../components/manage-resources/SummaryCard";

const ManageResources = () => (
    <Container fluid className="main-content-container px-4">
        {/* Page Header */}
        <Row noGutters className="page-header py-4">
            <PageTitle title="" subtitle="Manage Resources" className="text-sm-left mb-3"/>
        </Row>
        <Row>
            {/* Informaiton */}
            <Col lg="12" className="mb-4">
                <SummaryCard/>
            </Col>
        </Row>
        <Row>
            {/* Statistics */}
            <Col className="col-lg mb-4">
                <Statistics/>
            </Col>
        </Row>
    </Container>
);

ManageResources.propTypes = {};

ManageResources.defaultProps = {};

export default ManageResources;
