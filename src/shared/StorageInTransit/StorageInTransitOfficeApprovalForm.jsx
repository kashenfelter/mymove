import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export class StorageInTransitOfficeApprovalForm extends Component {
  render() {
    const { storageInTransitSchema } = this.props;
    return (
      <form
        data-cy="storage-in-transit-office-approval-form"
        onSubmit={this.props.handleSubmit(this.props.onSubmit)}
        className="storage-in-transit-form"
      >
        <fieldset key="sit-approval-information">
          <div className="editable-panel-column">
            <SwaggerField
              fieldName="authorized_start_date"
              swagger={storageInTransitSchema}
              title="Earliest authorized start date"
              required
            />
          </div>
          <div className="editable-panel-column">
            <SwaggerField
              className="sit-approval-field"
              fieldName="authorization_notes"
              title="Note"
              swagger={storageInTransitSchema}
            />
          </div>
        </fieldset>
      </form>
    );
  }
}

StorageInTransitOfficeApprovalForm.propTypes = {
  storageInTransitSchema: PropTypes.object.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

export const formName = 'storage_in_transit_office_approval_form';
StorageInTransitOfficeApprovalForm = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(StorageInTransitOfficeApprovalForm);

function mapStateToProps(state) {
  return {
    storageInTransitSchema: get(state, 'swaggerPublic.spec.definitions.StorageInTransit', {}),
  };
}

export default connect(mapStateToProps)(StorageInTransitOfficeApprovalForm);
