import logo from '../../images/logo/tides logo.png';
import PropTypes from 'prop-types';
import React from 'react';
import {Button, Form, FormGroup, Input, Label} from 'reactstrap';
import validate from "../../utils/validate";
import axios from "axios";
import classnames from "classnames";
import $ from "jquery";
import {Redirect} from "react-router-dom";

class AuthForm extends React.Component {

    state = {
        redirectToReferrer: false,
        formIsValid: false,
        formControls: {
            uname: {
                value: '',
                valid: false,
                validationRules: {
                    minLength: 2,
                    isRequired: true
                }
            },
            password: {
                value: '',
                valid: false,
                validationRules: {
                    minLength: 6,
                    isRequired: true
                }
            },
            repassword: {
                value: '',
                valid: false,
                validationRules: {
                    minLength: 6,
                    isRequired: true
                }
            },
        }
    };

    resetState = () => {
        const redirectToReferrer = false;
        const formIsValid = false;
        const formControls = {
            uname: {
                value: '',
                valid: false,
                validationRules: {
                    minLength: 2,
                    isRequired: true
                }
            },
            password: {
                value: '',
                valid: false,
                validationRules: {
                    minLength: 6,
                    isRequired: true
                }
            },
            repassword: {
                value: '',
                valid: false,
                validationRules: {
                    minLength: 6,
                    isRequired: true
                }
            },
        };
        this.setState({
            redirectToReferrer: redirectToReferrer,
            formControls: formControls,
            formIsValid: formIsValid
        });
    };

    get isLogin() {
        return this.props.authState === STATE_LOGIN;
    }

    get isSignup() {
        return this.props.authState === STATE_SIGNUP;
    }

    changeAuthState = authState => event => {
        event.preventDefault();

        this.props.onChangeAuthState(authState);
    };

    handleSubmit = event => {
        event.preventDefault();
        /**/
        $("#send").hide();
        // $("#action-feedback").html(`<img src = "../../images/theme/ajax-loader.gif">`);

        let all_filled = true;
        $("input[required]").each(function () {
            if ($(this).val() == "") {
                $(this).parent().css("border", "1px solid red");
                all_filled = all_filled && false;
            }
        });

        if (all_filled == false) {
            $("#action-feedback").html('');
            $("#send").show();
            return;
        }


        if ($("#password").val() != $("#repassword").val() && this.isSignup) {
            $("#action-feedback").html('Passwords do not match');
            $("#send").show();
            return;
        }

        if (!this.state.formControls.password.valid && !this.state.formControls.repassword.valid) {
            $("#action-feedback").html('Passwords too short!');
            $("#send").show();
            return;
        }

        const formData = {};
        for (let formElementId in this.state.formControls) {
            formData[formElementId] = this.state.formControls[formElementId].value;
        }


        // Handle Login Success or Fail or Register
        this.setState(() => ({
            redirectToReferrer: true
        }));
        ISAUTHENTICATED = true;

        // axios.get("/api/users/?format=json").then(({data}) => {
        //     console.log(data)
        // }).catch(err => console.log(err))
        // axios.post(`functions/register-handler`, formData)
        //     .then(res => {
        //         if (res.status == 200) {
        //             this.resetState(true);
        //             if (res.data == "1") {
        //                 $("#send").show();
        //                 window.location.href = "verification";
        //             } else {
        //                 $("#action-feedback").html(res.data);
        //                 $("#send").show();
        //             }
        //         } else {
        //             $("#action-feedback").html('Failed to create user. Please try again!');
        //             $("#send").show();
        //         }
        //     })
    };

    handleChange = event => {
        const name = event.target.name;
        let value = event.target.value;

        const updatedControls = {
            ...this.state.formControls
        };
        const updatedFormElement = {
            ...updatedControls[name]
        };

        updatedFormElement.value = value;
        updatedFormElement.valid = validate(value, updatedFormElement.validationRules);

        updatedControls[name] = updatedFormElement;

        let formIsValid = true;
        for (let inputIdentifier in updatedControls) {
            formIsValid = updatedControls[inputIdentifier].valid && formIsValid;
        }

        const passwordMatch = updatedControls["password"].value == updatedControls["repassword"].value;
        formIsValid = formIsValid && passwordMatch;

        $("input[required]").each(function () {
            $(this).parent().css("border", "none");
        });

        // $("#action-feedback").html('');
        if (!passwordMatch && updatedControls["password"].value && updatedControls["repassword"].value) {
            // $("#action-feedback").html('Passwords do not match');
        }

        this.setState({
            formControls: updatedControls,
            formIsValid: formIsValid
        });
    };

    renderButtonText() {
        const {buttonText} = this.props;

        if (!buttonText && this.isLogin) {
            return 'Login';
        }

        if (!buttonText && this.isSignup) {
            return 'Signup';
        }

        return buttonText;
    }

    render() {
        const {
            showLogo,
            usernameLabel,
            usernameInputProps,
            passwordLabel,
            passwordInputProps,
            confirmPasswordLabel,
            confirmPasswordInputProps,
            children,
            onLogoClick,
        } = this.props;

        const {redirectToReferrer} = this.state;

        if (redirectToReferrer === true) {
            return <Redirect to='/home'/>
        }


        return (
            <Form onSubmit={this.handleSubmit}>
                {showLogo && (
                    <div className="text-center pb-4">
                        <img
                            src={logo}
                            className=""
                            style={{width: 131, height: 88, cursor: 'pointer'}}
                            alt="logo"
                            onClick={onLogoClick}
                        />
                    </div>
                )}
                <FormGroup className={classnames({
                    focused: this.state.unameFocused
                })}>
                    {/*<Label for={usernameLabel}></Label>*/}
                    <Input
                        onChange={this.handleChange}
                        value={this.state.formControls.uname.value}
                        onFocus={e => this.setState({unameFocused: true})}
                        onBlur={e => this.setState({unameFocused: false})}
                        valid={this.state.formControls.uname.valid}
                        required
                        {...usernameInputProps} />
                </FormGroup>
                <FormGroup className={classnames({
                    focused: this.state.passwordFocused
                })}>
                    {/*<Label for={passwordLabel}></Label>*/}
                    <Input
                        onChange={this.handleChange}
                        value={this.state.formControls.password.value}
                        valid={this.state.formControls.password.valid}
                        onFocus={e => this.setState({passwordFocused: true})}
                        onBlur={e => this.setState({passwordFocused: false})}
                        required
                        {...passwordInputProps} />
                </FormGroup>
                {this.isSignup && (
                    <FormGroup className={classnames({
                        focused: this.state.repasswordFocused
                    })}>
                        {/*<Label for={confirmPasswordLabel}>{confirmPasswordLabel}</Label>*/}
                        <Input
                            onChange={this.handleChange}
                            value={this.state.formControls.repassword.value}
                            valid={this.state.formControls.repassword.valid}
                            onFocus={e => this.setState({repasswordFocused: true})}
                            onBlur={e => this.setState({repasswordFocused: false})}
                            required
                            {...confirmPasswordInputProps} />
                    </FormGroup>
                )}
                <Button
                    size="lg"
                    className="bg-gradient-theme-left border-0"
                    block
                    id="send"
                    onClick={this.handleSubmit}>
                    {this.renderButtonText()}
                </Button>

                <div className="text-center pt-1">
                    {this.isSignup ? (
                        <h6
                            className="mb-0 text-white text-regular"> Already registered?
                        </h6>
                    ) : (
                        <h6
                            className="mb-0 text-white text-regular"> No account?
                        </h6>
                    )}
                    <h6>
                        {this.isSignup ? (
                            <a href="#login" style={{color: 'black'}} onClick={this.changeAuthState(STATE_LOGIN)}>
                                sign in
                            </a>
                        ) : (
                            <a href="#signup" style={{color: 'black'}} onClick={this.changeAuthState(STATE_SIGNUP)}>
                                register account
                            </a>
                        )}
                    </h6>
                </div>
                <span id="action-feedback" className="text-danger"></span>
                {
                    children
                }
            </Form>
        );
    }
}

export const STATE_LOGIN = 'LOGIN';
export const STATE_SIGNUP = 'SIGNUP';
export var ISAUTHENTICATED = false;

AuthForm.propTypes = {
    authState: PropTypes.oneOf([STATE_LOGIN, STATE_SIGNUP]).isRequired,
    showLogo: PropTypes.bool,
    usernameLabel: PropTypes.string,
    usernameInputProps: PropTypes.object,
    passwordLabel: PropTypes.string,
    passwordInputProps: PropTypes.object,
    confirmPasswordLabel: PropTypes.string,
    confirmPasswordInputProps: PropTypes.object,
    onLogoClick: PropTypes.func,
};

AuthForm.defaultProps = {
    authState: 'LOGIN',
    showLogo: true,
    usernameLabel: 'Username',
    usernameInputProps: {
        name: 'uname',
        id: 'uname',
        type: 'text',
        placeholder: 'username',
    },
    passwordLabel: 'Password',
    passwordInputProps: {
        name: 'password',
        id: 'password',
        autoComplete: 'off',
        type: 'password',
        placeholder: 'password',
    },
    confirmPasswordLabel: 'Confirm Password',
    confirmPasswordInputProps: {
        name: 'repassword',
        id: 'repassword',
        autoComplete: 'off',
        type: 'password',
        placeholder: 'confirm password',
    },
    onLogoClick: () => {
    },
};

export default AuthForm;
