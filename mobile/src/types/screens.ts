import { ActionParse, ReactionParse, Workflow } from './aboutJson';

export type AuthParamList = {
  Login: undefined;
  Register: undefined;
  Home: undefined;
};

type OptionsProps = {
  isAction: boolean;
  setValues: (values: ActionOrReaction) => void;
};

type WorkflowDetailsProps = {
  workflow: Workflow;
};

export type AppStackList = {
  App: undefined;
  Options: OptionsProps;
  Home: undefined;
  Workflows: undefined;
  Service: undefined;
  Auth: undefined;
  Parameters: undefined;
  'Workflow Details': WorkflowDetailsProps;
  Setting: undefined;
};

// For Workflows Creation
export interface OptionValues {
  name: string;
  value: string | number;
  type: string;
}
export interface ActionOrReaction {
  id: number;
  name: string;
  description: string;
  options: OptionValues[];
}
