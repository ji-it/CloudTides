import React from "react";
import PropTypes from "prop-types";
import {Container, Row, Col} from "shards-react";

import PageTitle from "./../components/common/PageTitle";
import Policies from "../components/contribution/Policies";
import UsageResource from "../components/dashboard/UsageResource";
import ContributionCard from "../components/dashboard/ContributionCard";
import PowerCard from "../components/dashboard/PowerCard";
import WorkloadsCard from "../components/dashboard/WorkloadsCard";

const Contribution = () => (
    <Container fluid className="main-content-container px-4">
        {/* Page Header */}
        <Row noGutters className="page-header py-4">
            <PageTitle title="Contribution" subtitle="Manage Policies" className="text-sm-left mb-3"/>
        </Row>
        <Row>
            {/* Resources */}
            <Col className="col-lg mb-4">
                <Policies/>
            </Col>
        </Row>
    </Container>
);

Contribution.propTypes = {};

Contribution.defaultProps = {};

export default Contribution;
