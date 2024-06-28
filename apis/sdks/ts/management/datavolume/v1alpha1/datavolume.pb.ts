/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as VinkCommonCommon from "../../../common/common.pb"
import * as fm from "../../../fetch.pb"
import * as GoogleProtobufStruct from "../../../google/protobuf/struct.pb"
import * as GoogleProtobufTimestamp from "../../../google/protobuf/timestamp.pb"

type Absent<T, K extends keyof T> = { [k in Exclude<keyof T, K>]?: undefined };
type OneOf<T> =
  | { [k in keyof T]?: undefined }
  | (
    keyof T extends infer K ?
      (K extends string & keyof T ? { [k in K]: T[K] } & Absent<T, K>
        : never)
    : never);

export enum DataVolumeType {
  IMAGE = "IMAGE",
  ROOT = "ROOT",
  DATA = "DATA",
}

export enum OperatingSystemType {
  LINUX = "LINUX",
  WINDOWS = "WINDOWS",
  UBUNTU = "UBUNTU",
  CENTOS = "CENTOS",
  DEBIAN = "DEBIAN",
}

export enum OperatingSystemWindowsVersion {
  WINDOWS_10 = "WINDOWS_10",
  WINDOWS_11 = "WINDOWS_11",
}

export enum OperatingSystemUbuntuVersion {
  UBUNTU_18_04 = "UBUNTU_18_04",
  UBUNTU_20_04 = "UBUNTU_20_04",
  UBUNTU_22_04 = "UBUNTU_22_04",
}

export enum OperatingSystemCentOSVersion {
  CENTOS_7 = "CENTOS_7",
  CENTOS_8 = "CENTOS_8",
}

export enum OperatingSystemDebianVersion {
  DEBIAN_9 = "DEBIAN_9",
  DEBIAN_10 = "DEBIAN_10",
  DEBIAN_11 = "DEBIAN_11",
}

export type CreateDataVolumeRequest = {
  namespace?: string
  name?: string
  config?: DataVolumeConfig
}

export type DeleteDataVolumeRequest = {
  namespace?: string
  name?: string
}

export type DataVolumeConfigDataSourceBlank = {
}

export type DataVolumeConfigDataSourceUpload = {
}

export type DataVolumeConfigDataSourceHttp = {
  url?: string
  headers?: {[key: string]: string}
}

export type DataVolumeConfigDataSourceRegistry = {
  url?: string
}

export type DataVolumeConfigDataSourceS3 = {
  url?: string
}


/* vink modified */ export type BaseDataVolumeConfigDataSource = {
}

export type DataVolumeConfigDataSource = BaseDataVolumeConfigDataSource
  & OneOf<{ http: DataVolumeConfigDataSourceHttp; registry: DataVolumeConfigDataSourceRegistry; s3: DataVolumeConfigDataSourceS3; blank: DataVolumeConfigDataSourceBlank; upload: DataVolumeConfigDataSourceUpload }>

export type DataVolumeConfigBoundPVC = {
  storageClassName?: string
  capacity?: string
}


/* vink modified */ export type BaseDataVolumeConfigOperatingSystem = {
  type?: OperatingSystemType
}

export type DataVolumeConfigOperatingSystem = BaseDataVolumeConfigOperatingSystem
  & OneOf<{ windows: OperatingSystemWindowsVersion; ubuntu: OperatingSystemUbuntuVersion; centos: OperatingSystemCentOSVersion; debian: OperatingSystemDebianVersion }>

export type DataVolumeConfig = {
  dataVolumeType?: DataVolumeType
  operatingSystem?: DataVolumeConfigOperatingSystem
  dataSource?: DataVolumeConfigDataSource
  boundPvc?: DataVolumeConfigBoundPVC
}

export type DataVolume = {
  namespace?: string
  name?: string
  dataVolume?: GoogleProtobufStruct.Struct
  creationTimestamp?: GoogleProtobufTimestamp.Timestamp
}

export type DeleteDataVolumeResponse = {
}

export type ListDataVolumesRequest = {
  namespace?: string
  options?: VinkCommonCommon.ListOptions
}

export type ListDataVolumesResponse = {
  items?: DataVolume[]
  options?: VinkCommonCommon.ListOptions
}

export class DataVolumeManagement {
  static CreateDataVolume(req: CreateDataVolumeRequest, initReq?: fm.InitReq): Promise<DataVolume> {
    return fm.fetchReq<CreateDataVolumeRequest, DataVolume>(`/apis/vink.io/v1alpha1/namespaces/${req["namespace"]}/datavolumes/${req["name"]}`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static DeleteDataVolume(req: DeleteDataVolumeRequest, initReq?: fm.InitReq): Promise<DeleteDataVolumeResponse> {
    return fm.fetchReq<DeleteDataVolumeRequest, DeleteDataVolumeResponse>(`/apis/vink.io/v1alpha1/namespaces/${req["namespace"]}/datavolumes/${req["name"]}`, {...initReq, method: "DELETE"})
  }
  static ListDataVolumes(req: ListDataVolumesRequest, initReq?: fm.InitReq): Promise<ListDataVolumesResponse> {
    return fm.fetchReq<ListDataVolumesRequest, ListDataVolumesResponse>(`/apis/vink.io/v1alpha1/namespaces/${req["namespace"]}/datavolumes?${fm.renderURLSearchParams(req, ["namespace"])}`, {...initReq, method: "GET"})
  }
}