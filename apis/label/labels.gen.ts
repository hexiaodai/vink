
// GENERATED FILE -- DO NOT EDIT

export enum FeatureStatus {
    Alpha = "Alpha",
    Beta = "Beta",
    Stable = "Stable",
    Unknown = "Unknown"
}

export enum ResourceTypes {
    Unknown,
}

export class Instance {
    name: string;
    description: string;
    featureStatus: FeatureStatus;
    hidden: boolean;
    deprecated: boolean;
    resources: ResourceTypes[];

    constructor(
        name: string,
        description: string,
        featureStatus: FeatureStatus,
        hidden: boolean,
        deprecated: boolean,
        resources: ResourceTypes[]
    ) {
        this.name = name;
        this.description = description;
        this.featureStatus = featureStatus;
        this.hidden = hidden;
        this.deprecated = deprecated;
        this.resources = resources;
    }
}

export function allResourceLabels(): Instance[] {
    return [
    ];
}

export function allResourceTypes(): string[] {
    return [
    ];
}
