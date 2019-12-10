import React from "react";
import PropTypes from "prop-types";
import {
    Card,
    CardHeader,
    CardBody,
} from "shards-react";

class PowerCard extends React.Component {

    render() {
        const {title} = this.props;
        const nf1 = new Intl.NumberFormat('en-US', {
            style: 'percent',
        });
        const nf2 = new Intl.NumberFormat('en-US', {
            maximumFractionDigits: 2,
        });
        // const {percentage, wattage} = this.props.data;
        const {data} = this.props;
        return (
            <Card small className="h-100">
                <CardHeader className="border-bottom">
                    <h6 className="m-0">{title}</h6>
                </CardHeader>
                <CardBody className="d-flex py-0">
                    <div className="m-auto text-center">
                        <div className="mt-2" style={{fontSize: "1.2em", color: "#1B2376"}}>
                            <span style={{fontSize: "1.8em"}}>{nf2.format(data.power * 134)}</span> kWh
                        </div>
                        <div className="mb-3">
                            <span>{nf1.format(data.power)}</span> contributions
                        </div>
                    </div>
                </CardBody>
            </Card>
        );
    }
}

PowerCard.propTypes = {
    /**
     * The component's title.
     */
    title: PropTypes.string,
    data: PropTypes.object,
};

PowerCard.defaultProps = {
    title: "Power",
    data: {
        wattage: 1000,
        percentage: 0.6,
    }
};

export default PowerCard;
