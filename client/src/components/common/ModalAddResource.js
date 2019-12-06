import React from "react";
import {
    Button,
    Modal,
    ModalBody,
    ModalHeader,
    InputGroup,
    Form,
    FormInput,
    FormGroup,
    FormSelect,
    InputGroupAddon,
    InputGroupText
} from "shards-react";
import {Label, Input} from "reactstrap";
import classnames from "classnames";
import axios from 'axios';
import lodash from 'lodash'
import validate from "../../utils/validate";
import {Actions} from "../../flux";
import {element} from "prop-types";
import {devURL} from "../../utils/urls";
import request from "../../utils/request";
import auth from "../../utils/auth";

export default class ModalAddResource extends React.Component {

    state = {
        dcs: [],
        formIsValid: false,
        validateValid: false,
        formControls: {
            name: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            uname: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true,
                }
            },
            password: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            vmtype: {
                value: 'vSphere',
                valid: true,
                validationRules: {
                    isRequired: true
                }
            },
            host: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            datacenters: {
                value: "",
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
        }
    };

    resetState = (success) => {
        const formIsValid = false;
        const dcs = [];
        const formControls = {
            uname: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            name: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            password: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            vmtype: {
                value: "KVM",
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            host: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            datacenters: {
                value: "",
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
        };
        this.setState({
            formControls: formControls,
            validateValid: false,
            formIsValid: formIsValid,
            dcs: dcs
        });
        this.props.onExit("addModal", success);
    };

    handleChange = event => {
        const name = event.target.name;
        const value = event.target.value;

        const updatedControls = {
            ...this.state.formControls
        };
        const updatedFormElement = {
            ...updatedControls[name]
        };
        updatedFormElement.value = value;
        updatedFormElement.valid = validate(value, updatedFormElement.validationRules);

        updatedControls[name] = updatedFormElement;

        const {uname, password, host} = this.state.formControls;
        const validated = uname.valid && password.valid && host.valid;


        let formIsValid = validated;
        for (let inputIdentifier in updatedControls) {
            formIsValid = updatedControls[inputIdentifier].valid && formIsValid && this.state;
        }
        this.setState({
            formControls: updatedControls,
            formIsValid: formIsValid,
            validateValid: validated
        });
    };

    validateCredentials = event => {
        event.preventDefault();

        const {uname, password, host} = this.state.formControls;
        const formData = {"username": uname.value, "password": password.value, "host": host.value};
        const endpoint = '/api/resource/validate/';
        const requestURL = devURL + endpoint;
        request(requestURL, {method: 'POST', body: formData})
            .then((response) => {
                let formControls = this.state.formControls;
                if (response.results.length !== 0) {
                    formControls.datacenters.value = response.results[0];
                    formControls.datacenters.valid = true;
                }
                this.setState({
                    dcs: response.results,
                    formControls: formControls
                })
            }).catch((err) => {
            console.log(err);
        });
    };

    handleSubmit = event => {
        event.preventDefault();

        const formData = {};
        for (let formElementId in this.state.formControls) {
            formData[formElementId] = this.state.formControls[formElementId].value;
        }
        formData["polling_interval"] = 30;
        Actions.addResource(formData);
        this.resetState(true)
    };


    render() {

        return (
            <div>
                <Modal
                    open={this.props.toggleState}
                    size="md"
                    toggle={() => this.props.onExit("addModal", false)}
                >
                    <ModalHeader className="border-0 m-auto">
                        Add Resource
                    </ModalHeader>
                    <ModalBody className="border-top">
                        <Form role="form">
                            <FormGroup
                                className={classnames("mb-3", {
                                    focused: this.state.nameFocused
                                })}
                            >
                                <Label for="name">Name</Label>
                                <FormInput
                                    placeholder="Resource Name"
                                    type="text"
                                    name="name"
                                    onChange={this.handleChange}
                                    value={this.state.formControls.name.value}
                                    onFocus={e => this.setState({nameFocused: true})}
                                    onBlur={e => this.setState({nameFocused: false})}
                                    valid={this.state.formControls.name.valid}
                                />

                            </FormGroup>
                            <FormGroup
                                className={classnames("mb-3", {
                                    focused: this.state.vmtypeFocused
                                })}
                            >
                                <Label for="vmtype">VM Platform</Label>
                                <FormSelect
                                    name="vmtype"
                                    onChange={this.handleChange}
                                    value={this.state.formControls.vmtype.value}
                                    valid={this.state.formControls.vmtype.valid}
                                    onFocus={e => this.setState({vmtypeFocused: true})}
                                    onBlur={e => this.setState({vmtypeFocused: false})}
                                >
                                    <option value="vSphere">vSphere</option>
                                    <option value="KVM" disabled>KVM</option>
                                    <option value="Hyper-V" disabled>Hyper-V</option>
                                    <option value="XenServer" disabled>XenServer</option>
                                </FormSelect>
                            </FormGroup>
                            <FormGroup
                                className={classnames("mb-3", {
                                    focused: this.state.hostFocused
                                })}
                            >
                                <Label for="host">Host Address</Label>
                                <FormInput
                                    placeholder="FQDN or IP Address"
                                    type="text"
                                    name="host"
                                    onChange={this.handleChange}
                                    value={this.state.formControls.host.value}
                                    onFocus={e => this.setState({hostFocused: true})}
                                    onBlur={e => this.setState({hostFocused: false})}
                                    valid={this.state.formControls.host.valid}
                                />
                            </FormGroup>
                            <FormGroup>
                                <Label for="account">Account</Label>
                                <FormInput
                                    className={classnames({
                                        focused: this.state.unameFocused
                                    })}
                                    placeholder="Username"
                                    type="text"
                                    autoComplete="username"
                                    name="uname"
                                    onChange={this.handleChange}
                                    value={this.state.formControls.uname.value}
                                    valid={this.state.formControls.uname.valid}
                                    onFocus={e => this.setState({unameFocused: true})}
                                    onBlur={e => this.setState({unameFocused: false})}
                                />
                                <FormInput
                                    className={classnames({
                                        focused: this.state.passwordFocused
                                    })}
                                    placeholder="Password"
                                    type="password"
                                    name="password"
                                    autoComplete="current-password"
                                    onChange={this.handleChange}
                                    value={this.state.formControls.password.value}
                                    valid={this.state.formControls.password.valid}
                                    onFocus={e =>
                                        this.setState({passwordFocused: true})
                                    }
                                    onBlur={e =>
                                        this.setState({ppasswordFocused: false})
                                    }
                                />
                                <Button id="validate" className="my-2" disabled={!this.state.validateValid}
                                        onClick={this.validateCredentials} style={{"backgroundColor": "#00790e"}}
                                        type="button">
                                    Validate
                                </Button>
                            </FormGroup>
                            <FormGroup>
                                <Label for="datacenters">Select Datacenter</Label>
                                <Input
                                    className={classnames({
                                        focused: this.state.datacentersFocused
                                    })}
                                    type="select"
                                    name="datacenters"
                                    value={this.state.formControls.datacenters.value}
                                    onChange={this.handleChange}
                                    onFocus={e => this.setState({datacentersFocused: true})}
                                    onBlur={e => this.setState({datacentersFocused: false})}
                                >
                                    {
                                        this.state.dcs.map((item, index) => {
                                                return (<option value={item} key={index}>{item}</option>)
                                            }
                                        )}
                                </Input>
                            </FormGroup>
                            <div className="text-right">
                                <Button className="my-4 mr-2" theme="secondary" type="button"
                                        onClick={() => this.props.onExit("addModal", false)}
                                >Cancel</Button>
                                <Button id="send" className="my-4" disabled={!this.state.formIsValid}
                                        onClick={this.handleSubmit} color="primary"
                                        type="button">
                                    Save
                                </Button>
                            </div>
                        </Form>
                    </ModalBody>
                </Modal>
            </div>
        );
    }
}