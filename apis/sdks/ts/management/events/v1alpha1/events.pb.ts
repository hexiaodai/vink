/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../../fetch.pb"
export type Events = {
  data?: string
}

export type ListEventsRequest = {
}

export class EventsManagement {
  static List(req: ListEventsRequest, initReq?: fm.InitReq): Promise<Events> {
    return fm.fetchReq<ListEventsRequest, Events>(`/vink.kubevm.io.apis.management.events.v1alpha1.EventsManagement/List`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}