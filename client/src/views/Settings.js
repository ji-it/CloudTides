import React from "react";
import {Container, Row, Col} from "shards-react";

import PageTitle from "../components/common/PageTitle";
import UserDetails from "../components/settings/UserDetails";
import UserAccountDetails from "../components/settings/UserAccountDetails";
import Store from "../flux/store";
import {Actions} from "../flux";


class Settings extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            user: Store.getUserDetails(),
        };

        this.onChange = this.onChange.bind(this);
    }

    componentWillMount() {
        Store.addChangeListener(this.onChange);
    }

    componentWillUnmount() {
        Store.removeChangeListener(this.onChange);
    }

    componentDidMount() {
        Actions.getUserDetails();
    }

    onChange() {
        this.setState({
            ...this.state,
            user: Store.getUserDetails(),
        });
    }

    render() {
        const {user: userData} = this.state;

        return (
            <Container fluid className="main-content-container px-4">
                <Row noGutters className="page-header py-4">
                    <PageTitle title="User Profile" subtitle="Overview" md="12" className="ml-sm-auto mr-sm-auto"/>
                </Row>
                <Row>
                    <Col lg="4">
                        <UserDetails data={userData}/>
                    </Col>
                    <Col lg="8">
                        <UserAccountDetails data={userData}/>
                    </Col>
                </Row>
            </Container>
        );
    }
}

export default Settings;
