import React, { useState } from 'react';
import {
  View,
  Text,
  TextInput,
  StyleSheet,
  TouchableOpacity,
} from 'react-native';
import { globalStyles } from '../styles/global_style';
import { ActionOrReaction, AppStackList, Workflow } from '../types';
import { getWorkflows, sendWorkflows } from '../service';
import { NavigationProp, useNavigation } from '@react-navigation/native';

interface WorkflowCardProps {
  token: string;
  isBlackTheme?: boolean;
  serverIp: string;
  setWorkflows: (workflows: Workflow[]) => void;
  setRefresh: (refresh: boolean) => void;
}

export function WorkflowCard({
  isBlackTheme,
  serverIp,
  token,
  setWorkflows,
  setRefresh,
}: WorkflowCardProps) {
  const navigation = useNavigation<NavigationProp<AppStackList>>();

  const [workflowName, setWorkflowName] = useState('');
  const [action, setAction] = useState<ActionOrReaction>({
    id: 0,
    name: '',
    description: '',
    options: {},
  });
  const [reaction, setReaction] = useState<ActionOrReaction>({
    id: 0,
    name: '',
    description: '',
    options: {},
  });

  const handleSendWorkflow = async () => {
    if (token !== 'Error: token not found' && action && reaction) {
      await sendWorkflows(token, serverIp, {
        action_id: action.id,
        reaction_id: reaction.id,
        name: workflowName,
        action_option: action.options,
        reaction_option: reaction.options,
      });
      await getWorkflows(serverIp, token, setWorkflows);
      setAction({
        id: 0,
        name: '',
        description: '',
        options: {},
      });
      setReaction({
        id: 0,
        name: '',
        description: '',
        options: {},
      });
      setWorkflowName('');
      setRefresh(true);
    }
  };

  const isDisabled = action.name === '' || reaction.name === '';

  return (
    <View
      style={[
        styles.card,
        isBlackTheme ? globalStyles.secondaryLight : globalStyles.terciaryLight,
      ]}>
      <Text
        style={[
          isBlackTheme ? globalStyles.title : globalStyles.titleBlack,
          styles.title,
        ]}>
        Create a Workflow
      </Text>
      <TextInput
        style={[
          isBlackTheme ? globalStyles.input : globalStyles.inputBlack,
          styles.input,
        ]}
        placeholder="Workflow Name"
        placeholderTextColor={isBlackTheme ? '#0a0a0a' : 'f5f5f5'}
        value={workflowName}
        onChangeText={setWorkflowName}
      />
      <Text
        style={[
          isBlackTheme ? globalStyles.title : globalStyles.titleBlack,
          styles.title,
          { marginBottom: 20 },
        ]}>
        Select Action and Reaction
      </Text>
      <View style={styles.buttonContainer}>
        <TouchableOpacity
          style={[
            styles.ActionReactionButton,
            isBlackTheme
              ? globalStyles.secondaryDark
              : globalStyles.secondaryLight,
          ]}
          onPress={() => {
            navigation.navigate('Options', {
              isAction: true,
              setValues: setAction,
            });
          }}>
          {action.name === '' ? (
            <Text
              style={[
                isBlackTheme
                  ? globalStyles.textColorBlack
                  : globalStyles.textColor,
                globalStyles.textFormat,
              ]}
              numberOfLines={1}
              ellipsizeMode="tail">
              Add Action
            </Text>
          ) : (
            <Text
              style={[
                isBlackTheme
                  ? globalStyles.textColorBlack
                  : globalStyles.textColor,
                globalStyles.textFormat,
              ]}
              numberOfLines={1}
              ellipsizeMode="tail">
              {action.name}
            </Text>
          )}
        </TouchableOpacity>
        <TouchableOpacity
          style={[
            styles.ActionReactionButton,
            isBlackTheme
              ? globalStyles.secondaryDark
              : globalStyles.secondaryLight,
          ]}
          onPress={() => {
            navigation.navigate('Options', {
              isAction: false,
              setValues: setReaction,
            });
          }}>
          {reaction.name === '' ? (
            <Text
              style={[
                isBlackTheme
                  ? globalStyles.textColorBlack
                  : globalStyles.textColor,
                globalStyles.textFormat,
              ]}
              numberOfLines={1}
              ellipsizeMode="tail">
              Add Reaction
            </Text>
          ) : (
            <Text
              style={[
                isBlackTheme
                  ? globalStyles.textColorBlack
                  : globalStyles.textColor,
                globalStyles.textFormat,
              ]}
              numberOfLines={1}
              ellipsizeMode="tail">
              {reaction.name}
            </Text>
          )}
        </TouchableOpacity>
      </View>
      <TouchableOpacity
        disabled={isDisabled}
        style={[
          globalStyles.buttonFormat,
          { marginTop: 20 },
          isBlackTheme
            ? globalStyles.secondaryDark
            : globalStyles.secondaryLight,
          isDisabled && styles.disabledButton,
        ]}
        onPress={handleSendWorkflow}>
        <Text
          style={[
            isBlackTheme ? globalStyles.textColorBlack : globalStyles.textColor,
            globalStyles.textFormat,
          ]}>
          Send Workflow
        </Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    padding: 20,
    borderRadius: 10,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.8,
    shadowRadius: 2,
    elevation: 5,
    margin: 20,
  },
  title: {
    fontSize: 20,
  },
  input: {
    marginBottom: 20,
  },
  buttonContainer: {
    flexDirection: 'row',
    gap: 10,
  },
  ActionReactionButton: {
    width: 140,
    padding: 10,
    borderRadius: 20,
    justifyContent: 'center',
    alignItems: 'center',
    overflow: 'hidden',
  },
  disabledButton: {
    opacity: 0.5,
  },
});
