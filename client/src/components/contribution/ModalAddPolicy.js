import React from "react";
import {
    Button,
    Modal,
    ModalBody,
    ModalHeader,
    Collapse,
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
import {Actions} from "../../flux";

export default class ModalAddPolicy extends React.Component {

    state = {
        formIsValid: false,
        formControls: {
            name: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            username: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true,
                }
            },
            idle: {
                value: {
                    cpu: 0.30,
                    ram: 0.30,
                },
                isOpen: false,
                valid: true,
                validationRules: {}
            },
            threshold: {
                value: {
                    cpu: 0.60,
                    ram: 0.60,
                },
                isOpen: false,
                valid: true,
                validationRules: {}
            },
            accountType: {
                value: 'boinc',
                valid: true,
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
            isDestroy: {
                value: true,
                valid: true,
                validationRules: {
                    isRequired: true
                }
            },
            deployType: {
                value: 'VM',
                valid: true,
                validationRules: {
                    isRequired: true
                }
            },
            project: {
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
            name: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            username: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true,
                }
            },
            idle: {
                value: {
                    cpu: 0.30,
                    ram: 0.30,
                },
                isOpen: false,
                valid: true,
                validationRules: {}
            },
            threshold: {
                value: {
                    cpu: 0.60,
                    ram: 0.60,
                },
                isOpen: false,
                valid: true,
                validationRules: {}
            },
            accountType: {
                value: 'boinc',
                valid: true,
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
            isDestroy: {
                value: true,
                valid: true,
                validationRules: {
                    isRequired: true
                }
            },
            deployType: {
                value: 'VM',
                valid: true,
                validationRules: {
                    isRequired: true
                }
            },
            project: {
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

    toggleIdleCollapse = () => {
        const updatedControls = {
            ...this.state.formControls
        };
        updatedControls.idle.isOpen = !updatedControls.idle.isOpen;
        this.setState({formControls: updatedControls});
    };

    toggleThresholdCollapse = () => {
        const updatedControls = {
            ...this.state.formControls
        };
        updatedControls.threshold.isOpen = !updatedControls.threshold.isOpen;
        this.setState({formControls: updatedControls});
    };

    handleSlider = (name, type, values) => {
        const event = {target: {}};
        event.target.name = name;
        event.target.value = this.state.formControls[name].value;
        event.target.value[type] = Number(values[0])/100;
        this.handleChange(event)
    };

    handleIdleCPU = values => {
        this.handleSlider("idle", "cpu", values);
    };

    handleIdleRAM = values => {
        this.handleSlider("idle", "ram", values);
    };

    handleThreshCPU = values => {
        this.handleSlider("threshold", "cpu", values);
    };

    handleThreshRAM = values => {
        this.handleSlider("threshold", "ram", values);
    };

    handleSubmit = event => {
        event.preventDefault();

        const formData = {};
        for (let formElementId in this.state.formControls) {
            formData[formElementId] = this.state.formControls[formElementId].value;
        }
        Actions.addPolicy(formData);
        this.resetState(true)
    };


    render() {

        const acc_man_projects = [
            {name: "Science United", id: 1},
        ];

        const boinc_projects = [
            {name: "SETI@Home", id: 2},
        ];
        const IS_ACCMAN = this.state.formControls.accountType.value === "acc_manager";
        const projects = (IS_ACCMAN) ? acc_man_projects : boinc_projects;
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
                                className={classnames("mb-2", {
                                    focused: this.state.acctypeFocused
                                })}
                            >
                                <div><Label for="source">Choose Account Type</Label></div>
                                <FormRadio
                                    inline
                                    name="accountType"
                                    checked={this.state.formControls.accountType.value === "acc_manager"}
                                    onChange={this.handleChange}
                                    value="acc_manager"
                                >
                                    Account Manager
                                </FormRadio>
                                <FormRadio
                                    inline
                                    name="accountType"
                                    checked={this.state.formControls.accountType.value === "boinc"}
                                    value="boinc"
                                    onChange={this.handleChange}
                                >
                                    BOINC User
                                </FormRadio>
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
                                    name="username"
                                    onChange={this.handleChange}
                                    value={this.state.formControls.username.value}
                                    valid={this.state.formControls.username.valid}
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
                            </FormGroup>
                            <FormGroup
                                className={classnames("mb-3", {
                                    focused: this.state.projectFocused
                                })}
                            >
                                <Label for="project">Project</Label>
                                <FormSelect
                                    name="project"
                                    onChange={this.handleChange}
                                    value={this.state.formControls.project.value}
                                    valid={this.state.formControls.project.valid}
                                    onFocus={e => this.setState({projectFocused: true})}
                                    onBlur={e => this.setState({projectFocused: false})}
                                >
                                    <option value={""}></option>
                                    {projects.map((item, idx) => (
                                        <option key={idx} value={item.id}>{item.name}</option>
                                    ))}
                                </FormSelect>
                            </FormGroup>
                            <FormGroup
                                className={classnames("mb-3", {
                                    focused: this.state.deployFocused
                                })}
                            >
                                <Label for="deployType">Deploy Type</Label>
                                <FormSelect
                                    name="deployType"
                                    onChange={this.handleChange}
                                    value={this.state.formControls.deployType.value}
                                    valid={this.state.formControls.deployType.valid}
                                    onFocus={e => this.setState({deployFocused: true})}
                                    onBlur={e => this.setState({deployFocused: false})}
                                >
                                    <option value="VM">Virtual Machine</option>
                                    <option value="K8S">Kubernetes Cluster</option>
                                </FormSelect>
                            </FormGroup>
                            <FormGroup
                                className={classnames("mb-3", {
                                    focused: this.state.idleFocused
                                })}
                            >
                                <Label for="idle" onClick={this.toggleIdleCollapse}>Idle %</Label>
                                <Collapse open={this.state.formControls.idle.isOpen}>
                                    <ListGroupItem className="px-3">
                                        <div className="mb-2 pb-1">
                                            <div className="text-muted d-block mb-1">CPU</div>
                                            <Slider
                                                connect={[true, false]}
                                                start={[this.state.formControls.idle.value.cpu*100]}
                                                onChange={this.handleIdleCPU}
                                                range={{min: 0, max: 100}}
                                                tooltips
                                                name={"idle"}
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
                                                start={[this.state.formControls.idle.value.ram*100]}
                                                onChange={this.handleIdleRAM}
                                                range={{min: 0, max: 100}}
                                                tooltips
                                                name={"idle"}
                                                pips={{
                                                    mode: "positions",
                                                    values: [0, 25, 50, 75, 100],
                                                    stepped: true,
                                                    density: 5
                                                }}
                                            />
                                        </div>
                                    </ListGroupItem>
                                </Collapse>
                            </FormGroup>

                            <FormGroup
                                className={classnames("mb-3", {
                                    focused: this.state.idleFocused
                                })}
                            >
                                <Label for="threshold" onClick={this.toggleThresholdCollapse}>Threshold %</Label>
                                <Collapse open={this.state.formControls.threshold.isOpen}>
                                    <ListGroupItem className="px-3">
                                        <div className="mb-2 pb-1">
                                            <div className="text-muted d-block mb-1">CPU</div>
                                            <Slider
                                                connect={[true, false]}
                                                start={[this.state.formControls.threshold.value.cpu*100]}
                                                onChange={this.handleThreshCPU}
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
                                                start={[this.state.formControls.threshold.value.ram*100]}
                                                onChange={this.handleThreshRAM}
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
                                </Collapse>
                            </FormGroup>

                            <FormGroup
                                className={classnames("mb-3", {
                                    focused: this.state.isDestroyFocused
                                })}
                            >
                                <Label for="ip">Stop/Destroy when Busy</Label>
                                <div>
                                    <FormRadio
                                        inline
                                        name="isDestroy"
                                        checked={!this.state.formControls.isDestroy.value}
                                        onChange={this.handleChange}
                                    >
                                        Stop
                                    </FormRadio>
                                    <FormRadio
                                        inline
                                        name="isDestroy"
                                        checked={this.state.formControls.isDestroy.value}
                                        onChange={this.handleChange}
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