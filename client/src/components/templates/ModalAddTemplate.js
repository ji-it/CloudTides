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
    FormRadio,
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
import {offsetSign} from "recharts/lib/util/ChartUtils";

export default class ModalAddTemplate extends React.Component {

    state = {
        formIsValid: false,
        formControls: {
            source: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            compat: {
                value: '',
                valid: true,
                validationRules: {}
            },
            memsize: {
                value: 0,
                valid: true,
                validationRules: {}
            },
            os: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            uploaded_file: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            space: {
                value: 0,
                valid: true,
                validationRules: {}
            },
            vtemplate: {
                value: [],
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
        }
    };

    resetState = (success) => {
        const formIsValid = false;
        const formControls = {
            source: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            compat: {
                value: '',
                valid: true,
                validationRules: {}
            },
            memsize: {
                value: 0,
                valid: true,
                validationRules: {}
            },
            os: {
                value: '',
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            uploaded_file: {
                value: '',
                fileObj: {},
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
            space: {
                value: 0,
                valid: true,
                validationRules: {}
            },
            vtemplate: {
                value: [],
                valid: false,
                validationRules: {
                    isRequired: true
                }
            },
        };
        this.setState({
            formControls: formControls,
            formIsValid: formIsValid,
        });
        this.props.onExit("addModal", success);
    };

    handleChange = event => {
        const name = event.target.name;
        const value = event.target.value;
        const files = event.target.files;

        const updatedControls = {
            ...this.state.formControls
        };
        const updatedFormElement = {
            ...updatedControls[name]
        };
        updatedFormElement.value = value;

        if (name === 'uploaded_file' && files && files.length > 0)
            updatedFormElement.fileObj = files[0];

        updatedFormElement.valid = validate(value, updatedFormElement.validationRules);
        updatedControls[name] = updatedFormElement;

        let formIsValid = true;
        if (this.state.formControls.source.value === "upload") updatedControls.vtemplate.valid = true;

        for (let inputIdentifier in updatedControls) {
            formIsValid = updatedControls[inputIdentifier].valid && formIsValid;
        }
        this.setState({
            formControls: updatedControls,
            formIsValid: formIsValid,
        });
    };

    handleSubmit = event => {
        event.preventDefault();

        const formData = new FormData();
        const extra = {};
        for (let formElementId in this.state.formControls) {
            extra[formElementId] = this.state.formControls[formElementId].value;
            if (formElementId === "uploaded_file")
                formData.append('file', this.state.formControls[formElementId].fileObj);
        }
        formData.append('extra', JSON.stringify(extra));
        Actions.addTemplate(formData);
        this.resetState(true)
    };


    render() {
        const osTypes = ["", "Ubuntu Linux (64-bit)", "Linux"];
        const compatTypes = ["", "ESXI 6.5 and later", "ESXI 6.0"];
        return (
            <div>
                <Modal
                    open={this.props.toggleState}
                    size="md"
                    toggle={() => this.props.onExit("addModal", false)}
                >
                    <ModalHeader className="border-0 m-auto">
                        Add Template
                    </ModalHeader>
                    <ModalBody className="border-top">
                        <Form role="form">
                            <FormGroup
                                className={classnames("mb-3", {
                                    focused: this.state.sourceFocused
                                })}
                            >
                                <div><Label for="source">Choose Source</Label></div>
                                <FormRadio
                                    disabled
                                    inline
                                    name="source"
                                    checked={this.state.formControls.source.value === "datastore"}
                                    onChange={this.handleChange}
                                    value="datastore"
                                >
                                    From Datastore
                                </FormRadio>
                                <FormRadio
                                    inline
                                    name="source"
                                    checked={this.state.formControls.source.value === "upload"}
                                    value="upload"
                                    onChange={this.handleChange}
                                >
                                    Upload File
                                </FormRadio>
                            </FormGroup>
                            {
                                (this.state.formControls.source.value === "upload") ?
                                    (<div>
                                            <FormGroup
                                                className={classnames("mb-3", {
                                                    focused: this.state.fileFocused
                                                })}
                                            >
                                                <Label for="host">Select File</Label>
                                                <FormInput
                                                    placeholder="Choose file to upload"
                                                    type="file"
                                                    name="uploaded_file"
                                                    onChange={this.handleChange}
                                                    value={this.state.formControls.uploaded_file.value}
                                                    onFocus={e => this.setState({fileFocused: true})}
                                                    onBlur={e => this.setState({fileFocused: false})}
                                                    valid={this.state.formControls.uploaded_file.valid}
                                                />
                                            </FormGroup>
                                            <FormGroup
                                                className={classnames("mb-3", {
                                                    focused: this.state.osFocused
                                                })}
                                            >
                                                <Label for="os">Guest OS</Label>
                                                <FormSelect
                                                    name="os"
                                                    onChange={this.handleChange}
                                                    value={this.state.formControls.os.value}
                                                    valid={this.state.formControls.os.valid}
                                                    onFocus={e => this.setState({osFocused: true})}
                                                    onBlur={e => this.setState({osFocused: false})}
                                                >
                                                    {osTypes.map((item, index) => (
                                                        <option key={index} value={item}>{item}</option>
                                                    ))}
                                                </FormSelect>

                                            </FormGroup>
                                            <FormGroup
                                                className={classnames("mb-3", {
                                                    focused: this.state.compatFocused
                                                })}
                                            >
                                                <Label for="compat">Compatibility (optional)</Label>
                                                <FormSelect
                                                    name="compat"
                                                    onChange={this.handleChange}
                                                    value={this.state.formControls.compat.value}
                                                    // valid={this.state.formControls.compat.valid}
                                                    onFocus={e => this.setState({compatFocused: true})}
                                                    onBlur={e => this.setState({compatFocused: false})}
                                                >
                                                    {compatTypes.map((item, index) => (
                                                        <option key={index} value={item}>{item}</option>
                                                    ))}
                                                </FormSelect>
                                            </FormGroup>
                                            <FormGroup>
                                                <Label for="memsize">Memory Size - GB (optional)</Label>
                                                <FormInput
                                                    className={classnames({
                                                        focused: this.state.memsizeFocused
                                                    })}
                                                    type="number"
                                                    name="memsize"
                                                    value={this.state.formControls.memsize.value}
                                                    onChange={this.handleChange}
                                                    onFocus={e => this.setState({memsizeFocused: true})}
                                                    onBlur={e => this.setState({memsizeFocused: false})}
                                                />
                                            </FormGroup>
                                            <FormGroup>
                                                <Label for="space">Provisioned Space - GB (optional) </Label>
                                                <FormInput
                                                    className={classnames({
                                                        focused: this.state.spaceFocused
                                                    })}
                                                    type="number"
                                                    name="space"
                                                    value={this.state.formControls.space.value}
                                                    onChange={this.handleChange}
                                                    onFocus={e => this.setState({spaceFocused: true})}
                                                    onBlur={e => this.setState({spaceFocused: false})}
                                                />
                                            </FormGroup>
                                        </div>
                                    ) : (
                                        (this.state.formControls.source.value === "datastore") ?
                                            <div>
                                                <FormGroup
                                                    className={classnames("mb-3", {
                                                        focused: this.state.vtemplateFocused
                                                    })}
                                                >
                                                    <Label for="vtemplate">Choose template from datastore</Label>
                                                    <FormSelect
                                                        multiple
                                                        name="vtemplate"
                                                        onChange={this.handleChange}
                                                        value={this.state.formControls.vtemplate.value}
                                                        // valid={this.state.formControls.compat.valid}
                                                        onFocus={e => this.setState({vtemplateFocused: true})}
                                                        onBlur={e => this.setState({vtemplateFocused: false})}
                                                    >
                                                        {compatTypes.map((item, index) => (
                                                            <option key={index} value={item}>{item}</option>
                                                        ))}
                                                    </FormSelect>
                                                </FormGroup>
                                            </div>
                                            : <></>
                                    )
                            }
                            <div className="text-right">
                                <Button className="my-4 mr-2" theme="secondary" type="button"
                                        onClick={() => this.props.onExit("addModal", false)}
                                >Cancel</Button>
                                <Button id="send" className="my-4" disabled={!this.state.formIsValid}
                                        onClick={this.handleSubmit} color="primary"
                                        type="button">
                                    Upload
                                </Button>
                            </div>
                        </Form>
                    </ModalBody>
                </Modal>
            </div>
        );
    }
}