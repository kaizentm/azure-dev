import yaml from "yamljs";
import { IOptions } from "glob";
import path from "path";
import os from "os";
import fs from "fs/promises";
import { existsSync } from "fs";
import ansiEscapes from "ansi-escapes";
import chalk from "chalk";
import { cleanDirectoryPath, copyFile, createRepoUrlFromRemote, ensureDirectoryPath, getGlobFiles, getRepoPropsFromRemote, isStringNullOrEmpty, RepoProps, writeHeader } from "../common/util";
import { AssetRule, GitRemote, RepomanCommand, RepomanCommandOptions, RepoManifest } from "../models";
import { GitRepo } from "../tools/git";


export interface CleanCommandOptions extends RepomanCommandOptions {
    templateFile: string
    branch: string
    source: string
    output: string
    failOnCleanError?: boolean
}


export class CleanCommand implements RepomanCommand {
    private templateFile: string;
    private manifest: RepoManifest;
    private sourcePath: string;
    private outputPath: string;
    constructor(private options: CleanCommandOptions) {
        this.sourcePath = path.resolve(path.normalize(options.source))
        this.templateFile = path.join(this.sourcePath, options.templateFile);

        try {
            this.manifest = yaml.load(this.templateFile);
            this.outputPath = path.resolve(path.normalize(options.output))
        }
        catch (err) {
            console.error(chalk.red(`Repo template manifest not found at '${this.templateFile}'`));
            throw err;
        }
    }

    public execute = async () => {
        writeHeader(`Clean Command`);
    
        if(!this.validRemotes())
          return;

        console.info(chalk.white(`Repo clean started...`));

        this.manifest.repo.remotes.forEach(async remote => {
            try {
                await this.deleteRemoteBranch(remote);
                console.info(chalk.cyan(`Branch ${this.options.branch} has been deleted from remote ${remote.url}.`));
            }
            catch (err){
                console.error(chalk.red(err));
                if (this.options.failOnCleanError) {
                    throw err;
                }
            }
        });
        
        console.info(chalk.white(`Repo clean completed.`));
        console.info();
    }
    
    private deleteRemoteBranch = async (remote: GitRemote) => {
        const targetBranch: string = this.options.branch;
        const repoName: string = this.manifest.metadata.name;

        await ensureDirectoryPath(this.outputPath);
        await cleanDirectoryPath(this.outputPath);

        const repo = new GitRepo(this.outputPath);
        await repo.clone(repoName,remote.url);

        if(!await repo.remoteBranchExists(remote.url,targetBranch)){
            const message = `Error deleting remote branch ${targetBranch}. Branch does not exist on remote ${remote.url}`;
            throw message;
        }

        await repo.deleteRemoteBranch(repoName,targetBranch);
    }
    private validRemotes = (): Boolean => {
        if (!this.manifest.repo.remotes || this.manifest.repo.remotes.length === 0) {
            console.warn(chalk.yellowBright("Remotes manifest is missing 'remotes' configuration and is unable to push changes"));
            return false;
        }
        return true;
    }
}
