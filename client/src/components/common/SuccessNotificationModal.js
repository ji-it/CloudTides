import React from "react";
import {
    Button,
    Modal,
} from "shards-react";


class SuccessNotificationModal extends React.Component {

    componentDidMount() {
        this.props.onRef(this)
    }

    componentWillUnmount() {
        this.props.onRef(undefined)
    }

    state = {};
    toggleModal = state => {
        this.setState({
            [state]: !this.state[state]
        });
    };

    render() {
        return (
            <>
                <Modal
                    className="modal-dialog-centered modal-danger"
                    contentClassName="bg-gradient-success"
                    isOpen={this.state.notificationModal}
                    toggle={() => this.toggleModal("notificationModal")}
                >
                    <div className="modal-header">
                        <h6 className="modal-title" id="modal-title-notification">
                            Success!
                        </h6>
                        <button
                            aria-label="Close"
                            className="close"
                            data-dismiss="modal"
                            type="button"
                            onClick={() => this.toggleModal("notificationModal")}
                        >
                            <span aria-hidden={true}>×</span>
                        </button>
                    </div>
                    <div className="modal-body">
                        <div className="py-3 text-center">
                            <i className="ni ni-satisfied ni-3x"/>
                            <h4 className="heading mt-4">Email Successfully Sent!</h4>
                            <p>
                                Thank you for reaching out to us. We will contact you in no time!
                            </p>
                        </div>
                    </div>
                    <div className="modal-footer text-center">
                        <Button
                            className="text-white text"
                            color="link"
                            data-dismiss="modal"
                            type="button"
                            onClick={() => this.toggleModal("notificationModal")}
                            style={{margin: "0 auto"}}
                        >
                            Close
                        </Button>
                    </div>
                </Modal>
            </>
        )
    }
}

export default SuccessNotificationModal;