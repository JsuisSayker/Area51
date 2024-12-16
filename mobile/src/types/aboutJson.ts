export type Action = {
  action_id: number;
  name: string;
  description: string;
};

export type Reaction = {
  reaction_id: number;
  name: string;
  description: string;
};

export type Service = {
  name: string;
  actions: Action[];
  reactions: Reaction[];
};

export type Workflow = {
  name: string;
  action_id: number;
  reaction_id: number;
  is_active: boolean;
  created_at: string;
};
export type Server = {
  current_time: string;
  services: Service[];
  workflows: Workflow[];
};

export type Client = {
  host: string;
};

export type AboutJson = {
  client: Client;
  server: Server;
};

export type PullRequestComment = {
  body: string;
  pull_request_url: string;
};
