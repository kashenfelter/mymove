// eslint-disable-next-line no-unused-vars
import React from 'react';
import PropTypes from 'prop-types';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import 'scenes/Shipments/ShipmentCards.css';

const ShipmentCards = ({ shipments }) => {
  if (!shipments) return <LoadingPlaceholder />;
  if (shipments.length === 0)
    return <h2> There is no feedback at the moment! </h2>;
  return (
    <div className="shipment-cards">
      {shipments.map(shipment => (
        <div key={shipment.id} className="shipment-card">
          {shipment.name}
        </div>
      ))}
    </div>
  );
};

ShipmentCards.propTypes = {
  shipments: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      name: PropTypes.string.isRequired,
      traffic_distribution_list: PropTypes.string,
    }),
  ),
};

export default ShipmentCards;
