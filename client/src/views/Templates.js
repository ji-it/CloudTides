import React from "react";
import PropTypes from "prop-types";
import {Container, Row, Col} from "shards-react";

import PageTitle from "./../components/common/PageTitle";
import TemplatesTable from "../components/templates/TemplatesTable";

const Templates = () => (
    <Container fluid className="main-content-container px-4">
        {/* Page Header */}
        <Row noGutters className="page-header py-4">
            <PageTitle title="Templates" subtitle="Manage Templates" className="text-sm-left mb-3"/>
        </Row>
        <Row>
            {/* Resources */}
            <Col className="col-lg mb-4">
                <TemplatesTable/>
            </Col>
        </Row>
    </Container>
);

Templates.propTypes = {};

Templates.defaultProps = {};

export default Templates;
