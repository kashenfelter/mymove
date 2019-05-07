import { denormalize } from 'normalizr';
import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient, getClient } from 'shared/Swagger/api';
import { uploads } from '../schema';
import { ADD_ENTITIES } from '../actions';

export const STATE_KEY = 'uploads';

export default function reducer(state = {}, action) {
  switch (action.type) {
    case ADD_ENTITIES:
      return {
        ...state,
        ...action.payload.uploads,
      };

    default:
      return state;
  }
}

export const createShipmentDocumentLabel = 'Uploads.createShipmentDocument';
export const updateUploadLabel = 'uploads.updateUpload';

export function createShipmentDocument(shipmentId, createGenericMoveDocument, label = createShipmentDocumentLabel) {
  return swaggerRequest(
    getPublicClient,
    'move_docs.createGenericMoveDocument',
    { shipmentId, createGenericMoveDocument },
    { label },
  );
}

export const selectUpload = (state, id) => {
  return denormalize([id], uploads, state.entities)[0];
};

export function updateUpload(uploadId, params, label = updateUploadLabel) {
  return swaggerRequest(getClient, label, { uploadId, params }, { label });
}
