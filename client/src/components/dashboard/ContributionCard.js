import React from "react";
import PropTypes from "prop-types";
import {
    Card,
    CardHeader,
    CardBody,
} from "shards-react";

class ContributionCard extends React.Component {

    render() {
        const {title} = this.props;
        const nf = new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: "USD",
            minimumFractionDigits: 0,
            maximumFractionDigits: 0,
        });
        const {dayValue, monthValue} = this.props.data;
        return (
            <Card small className="h-100">
                <CardHeader className="border-bottom">
                    <h6 className="m-0">{title}</h6>
                </CardHeader>
                <CardBody className="d-flex py-0">
                    <div className="m-auto text-center">
                        <div className="mt-2" style={{fontSize: "1.2em", color: "#1B2376"}}>
                            <span style={{fontSize: "1.8em"}}>{nf.format(dayValue)}</span> /day
                        </div>
                        <div className="mb-3">
                            <span>{nf.format(monthValue)}</span> /month
                        </div>
                    </div>
                </CardBody>
            </Card>
        );
    }
}

ContributionCard.propTypes = {
    /**
     * The component's title.
     */
    title: PropTypes.string,
    data: PropTypes.object,
};

ContributionCard.defaultProps = {
    title: "Contribution",
    data: {
        dayValue: 1000,
        monthValue: 1000000,
    }
};

export default ContributionCard;
