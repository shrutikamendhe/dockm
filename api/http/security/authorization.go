package security

import "github.com/shrutikamendhe/dockm/api"

// AuthorizedResourceControlDeletion ensure that the user can delete a resource control object.
// A non-administrator user cannot delete a resource control where:
// * the AdministratorsOnly flag is set
// * he is not one of the users in the user accesses
// * he is not a member of any team within the team accesses
func AuthorizedResourceControlDeletion(resourceControl *dockm.ResourceControl, context *RestrictedRequestContext) bool {
	if context.IsAdmin {
		return true
	}

	if resourceControl.AdministratorsOnly {
		return false
	}

	userAccessesCount := len(resourceControl.UserAccesses)
	teamAccessesCount := len(resourceControl.TeamAccesses)

	if teamAccessesCount > 0 {
		for _, access := range resourceControl.TeamAccesses {
			for _, membership := range context.UserMemberships {
				if membership.TeamID == access.TeamID {
					return true
				}
			}
		}
	}

	if userAccessesCount > 0 {
		for _, access := range resourceControl.UserAccesses {
			if access.UserID == context.UserID {
				return true
			}
		}
	}

	return false
}

// AuthorizedResourceControlUpdate ensure that the user can update a resource control object.
// It reuses the creation restrictions and adds extra checks.
// A non-administrator user cannot update a resource control where:
// * he wants to put one or more user in the user accesses
func AuthorizedResourceControlUpdate(resourceControl *dockm.ResourceControl, context *RestrictedRequestContext) bool {
	userAccessesCount := len(resourceControl.UserAccesses)
	if !context.IsAdmin && userAccessesCount > 0 {
		return false
	}

	return AuthorizedResourceControlCreation(resourceControl, context)
}

// AuthorizedResourceControlCreation ensure that the user can create a resource control object.
// A non-administrator user cannot create a resource control where:
// * the AdministratorsOnly flag is set
// * he wants to add more than one user in the user accesses
// * he wants to add a team he is not a member of
func AuthorizedResourceControlCreation(resourceControl *dockm.ResourceControl, context *RestrictedRequestContext) bool {
	if context.IsAdmin {
		return true
	}

	if resourceControl.AdministratorsOnly {
		return false
	}

	userAccessesCount := len(resourceControl.UserAccesses)
	teamAccessesCount := len(resourceControl.TeamAccesses)
	if userAccessesCount > 1 || (userAccessesCount == 1 && teamAccessesCount == 1) {
		return false
	}

	if userAccessesCount == 1 {
		access := resourceControl.UserAccesses[0]
		if access.UserID == context.UserID {
			return true
		}
	}

	if teamAccessesCount > 0 {
		for _, access := range resourceControl.TeamAccesses {
			isMember := false
			for _, membership := range context.UserMemberships {
				if membership.TeamID == access.TeamID {
					isMember = true
				}
			}
			if !isMember {
				return false
			}
		}
	}

	return true
}

// AuthorizedTeamManagement ensure that access to the management of the specified team is granted.
// It will check if the user is either administrator or leader of that team.
func AuthorizedTeamManagement(teamID dockm.TeamID, context *RestrictedRequestContext) bool {
	if context.IsAdmin {
		return true
	}

	for _, membership := range context.UserMemberships {
		if membership.TeamID == teamID && membership.Role == dockm.TeamLeader {
			return true
		}
	}

	return false
}

// AuthorizedUserManagement ensure that access to the management of the specified user is granted.
// It will check if the user is either administrator or the owner of the user account.
func AuthorizedUserManagement(userID dockm.UserID, context *RestrictedRequestContext) bool {
	if context.IsAdmin || context.UserID == userID {
		return true
	}
	return false
}
