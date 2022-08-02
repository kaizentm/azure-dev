export interface RepoManifest {
    metadata: {
        type: string
        name: string
        description: string
    }
    repo: {
        includeProjectAssets: boolean
        remotes: GitRemote[]
        iac: IAC[]
        assets: AssetRule[]
        rewrite?: {
            patterns: string[]
        }
    }
}

export interface GitRemote {
    name: string
    url: string
    branch?: string
}

export interface AssetRule {
    from: string
    to: string
    patterns?: string[]
    ignore?: string[]
}

export interface RepomanCommandOptions {
    debug: true
    [key: string]: any
}

export interface RepomanCommand {
    execute(): Promise<void>
}

export interface IAC {
    name: string
    path : string
    updateRemoteUrl: boolean
    remoteSlug?: string
}

export enum IACTools{
    Bicep = "bicep",
    Terraform = "terraform"
}