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
    ListGroupItem,
    Slider,
    FormRadio,
    InputGroupAddon,
    InputGroupText
} from "shards-react";
import {Label, Input} from "reactstrap";
import classnames from "classnames";
import axios from 'axios';
import validate from "../../utils/validate";

export default class ModalAddPolicy extends React.Component {

    state = {
        formIsValid: false,
        formControls: {
            name: {
                value: '',
                valid: false,
                validationRules: {
                    minLength: 4,
                    isRequired: true
                }
            },
            uname: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true,
                    minLength: 2,
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
            vmtype: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            ip: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            }
        }
    };

    resetState = (success) => {
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
            name: {
                value: '',
                valid: false,
                validationRules: {
                    minLength: 4,
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
            vmtype: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            ip: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            }
        };
        this.setState({
            formControls: formControls,
            formIsValid: formIsValid
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

        let formIsValid = true;
        for (let inputIdentifier in updatedControls) {
            formIsValid = updatedControls[inputIdentifier].valid && formIsValid;
        }
        this.setState({
            formControls: updatedControls,
            formIsValid: formIsValid
        });
    };

    handleSubmit = event => {
        event.preventDefault();

        const formData = {};
        for (let formElementId in this.state.formControls) {
            formData[formElementId] = this.state.formControls[formElementId].value;
        }

        axios.post(`additional-functions/contact-form-handler`, formData)
            .then(res => {
                if (res.status == 200) {
                    this.resetState(true);
                } else {
                    console.log(res.data);
                }
            })
    };


    render() {

        const dcs = [
            {name: "NY Datacenter", id: 214},
            {name: "LA Datacenter", id: 213},
        ];

        return (
            <div>
                <Modal
                    open={this.props.toggleState}
                    size="md"
                    toggle={() => this.props.onExit("addModal", false)}
                >
                    <ModalHeader className="border-0 m-auto">
                        Add Policy
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
                                    placeholder="Policy Name"
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
                                    focused: this.state.emailFocused
                                })}
                            >
                                <Label for="email">Project</Label>
                                <FormSelect
                                    name="vmtype"
                                    onChange={this.handleChange}
                                    // value={this.state.formControls.vmtype.value}
                                    // valid={this.state.formControls.vmtype.valid}
                                    onFocus={e => this.setState({vmtypeFocused: true})}
                                    onBlur={e => this.setState({vmtypeFocused: false})}
                                >
                                    <option value="vSphere">SETI@Home</option>
                                    <option value="KVM" disabled>Cloud@Home</option>
                                </FormSelect>
                            </FormGroup>
                            <FormGroup
                                className={classnames("mb-3", {
                                    focused: this.state.emailFocused
                                })}
                            >
                                <Label for="email">Deploy Type</Label>
                                <FormSelect
                                    name="vmtype"
                                    onChange={this.handleChange}
                                    // value={this.state.formControls.vmtype.value}
                                    // valid={this.state.formControls.vmtype.valid}
                                    onFocus={e => this.setState({vmtypeFocused: true})}
                                    onBlur={e => this.setState({vmtypeFocused: false})}
                                >
                                    <option value="vSphere">Kubernetes Cluster</option>
                                    <option value="KVM" disabled>Virtual Machine</option>
                                </FormSelect>
                            </FormGroup>
                            <FormGroup
                                className={classnames("mb-3", {
                                    focused: this.state.emailFocused
                                })}
                            >
                                <Label for="email">Idle %</Label>
                                <ListGroupItem className="px-3">
                                    <div className="mb-2 pb-1">
                                        <div className="text-muted d-block mb-1">CPU</div>
                                        <Slider
                                            connect={[true, false]}
                                            start={[85]}
                                            range={{min: 0, max: 100}}
                                            tooltips
                                            pips={{
                                                mode: "positions",
                                                values: [0, 25, 50, 75, 100],
                                                stepped: true,
                                                density: 5
                                            }}
                                        />
                                    </div>
                                    <div className="mb-2 pb-1">
                                        <div className="text-muted d-block mb-1">Disk</div>
                                        <Slider
                                            connect={[true, false]}
                                            start={[50]}
                                            range={{min: 0, max: 100}}
                                            tooltips
                                            pips={{
                                                mode: "positions",
                                                values: [0, 25, 50, 75, 100],
                                                stepped: true,
                                                density: 5
                                            }}
                                        />
                                    </div>
                                    <div className="mb-2 pb-1">
                                        <div className="text-muted d-block mb-1">Memory</div>
                                        <Slider
                                            connect={[true, false]}
                                            start={[20]}
                                            range={{min: 0, max: 100}}
                                            tooltips
                                            pips={{
                                                mode: "positions",
                                                values: [0, 25, 50, 75, 100],
                                                stepped: true,
                                                density: 5
                                            }}
                                        />
                                    </div>
                                </ListGroupItem>
                            </FormGroup>
                            <FormGroup
                                className={classnames("mb-3", {
                                    focused: this.state.ipFocused
                                })}
                            >
                                <Label for="ip">Stop/Destroy when Busy</Label>
                                <div>
                                    <FormRadio
                                        inline
                                        name="sport"
                                        checked={true}
                                        // checked={this.state.selectedSport === "basketball"}
                                        // onChange={() => {
                                        //     this.changeSport("basketball");
                                        //}
                                    >
                                        Stop
                                    </FormRadio>
                                    <FormRadio
                                        inline
                                        name="sport"
                                        // checked={this.state.selectedSport === "basketball"}
                                        // onChange={() => {
                                        //     this.changeSport("basketball");
                                        //}
                                    >
                                        Destroy
                                    </FormRadio>
                                </div>
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