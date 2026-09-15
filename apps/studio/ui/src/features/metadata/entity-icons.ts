import type { Component } from 'vue'
import {
  Activity,
  ArrowUpRight,
  Box,
  Boxes,
  BriefcaseBusiness,
  CalendarClock,
  CircleDollarSign,
  Clock,
  Columns3,
  FileText,
  Flag,
  GitPullRequestArrow,
  Hash,
  House,
  KeyRound,
  Languages,
  LayoutDashboard,
  ListChecks,
  ListTree,
  Package,
  Settings2,
  Shield,
  ShieldCheck,
  User,
  UsersRound,
} from '@lucide/vue'

const entityIconRegistry: Record<string, Component> = {
  activity: Activity,
  'arrow-up-right': ArrowUpRight,
  box: Box,
  boxes: Boxes,
  'briefcase-business': BriefcaseBusiness,
  'calendar-clock': CalendarClock,
  'circle-dollar-sign': CircleDollarSign,
  clock: Clock,
  'columns-3': Columns3,
  'file-text': FileText,
  flag: Flag,
  'git-pull-request-arrow': GitPullRequestArrow,
  hash: Hash,
  house: House,
  'key-round': KeyRound,
  languages: Languages,
  'layout-dashboard': LayoutDashboard,
  'list-checks': ListChecks,
  'list-tree': ListTree,
  package: Package,
  'settings-2': Settings2,
  shield: Shield,
  'shield-check': ShieldCheck,
  user: User,
  'users-round': UsersRound,
}

export function iconForEntity(icon?: string): Component {
  const key = icon ? toEntityIconKey(icon) : ''
  return key ? entityIconRegistry[key] ?? Box : Box
}

function toEntityIconKey(value: string): string {
  return value
    .trim()
    .replace(/[_\s]+/g, '-')
    .replace(/([a-z])([A-Z0-9])/g, '$1-$2')
    .replace(/([0-9])([A-Z])/g, '$1-$2')
    .toLowerCase()
}
