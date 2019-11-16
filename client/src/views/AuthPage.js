import AuthForm, {STATE_LOGIN} from '../components/auth/AuthForn';
import React from 'react';
import {Card, Col, Row} from 'reactstrap';

class AuthPage extends React.Component {
    handleAuthState = authState => {
        if (authState === STATE_LOGIN) {
            this.props.history.push('/login');
        } else {
            this.props.history.push('/signup');
        }
    };


    handleLogoClick = () => {
        this.props.history.push('/');
    };

    render() {
        return (
            <Row
                style={{
                    height: '100vh',
                    justifyContent: 'center',
                    alignItems: 'center',
                    marginRight: "0px"
                }}>
                <Col md={6} lg={4}>
                    <div>
                        <AuthForm
                            authState={this.props.authState}
                            onChangeAuthState={this.handleAuthState}
                            onLogoClick={this.handleLogoClick}
                        />
                    </div>
                </Col>
            </Row>
        );
    }
}

export default AuthPage;
