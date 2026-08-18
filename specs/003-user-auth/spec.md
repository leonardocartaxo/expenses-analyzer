# Feature Specification: User Auth

**Feature Branch**: `003-user-auth`

**Created**: 2026-08-18

**Status**: Draft

**Input**: Self-serve sign-up, email/password sign-in, and password reset by email per docs/PRODUCT.md. No SSO. No email verification.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Sign up (Priority: P1)

A new person creates an account with email, password, name, and address (free text). They are not an application admin and belong to no organization yet.

**Why this priority**: No other product flow starts without an account.

**Independent Test**: Sign up with valid fields; account exists; can be used to sign in; isAppAdmin is false.

**Acceptance Scenarios**:

1. **Given** no account for an email, **When** the person submits sign-up with email, password, name, and address, **Then** an account is created and they can sign in with that email and password.
2. **Given** sign-up, **When** email is already used, **Then** the product rejects the attempt in plain language and does not overwrite the existing account.
3. **Given** sign-up, **When** name or address is empty, **Then** the product rejects the attempt.
4. **Given** a newly signed-up account, **When** someone inspects platform role, **Then** the user is not an application admin.

---

### User Story 2 - Sign in and sign out (Priority: P1)

An existing user signs in with email and password and reaches an authenticated session. They can sign out.

**Why this priority**: Required to create or join organizations later.

**Independent Test**: Sign in with correct credentials; signed-out user cannot access authenticated pages/APIs.

**Acceptance Scenarios**:

1. **Given** a signed-up user, **When** they submit the correct email and password, **Then** they are authenticated.
2. **Given** a signed-up user, **When** they submit a wrong password, **Then** they are not authenticated and see a plain-language error that does not reveal whether the email exists if that would leak accounts (same generic failure is acceptable).
3. **Given** an authenticated user, **When** they sign out, **Then** subsequent authenticated requests fail until they sign in again.

---

### User Story 3 - Stay signed in for the session (Priority: P2)

After sign-in, the user remains authenticated across page loads until they sign out or the session ends.

**Why this priority**: Needed for later org/bill work.

**Independent Test**: Sign in, reload, still authenticated; after sign out, not authenticated.

**Acceptance Scenarios**:

1. **Given** a signed-in user, **When** they reload the app, **Then** they are still authenticated without typing the password again.
2. **Given** no valid session, **When** they request an authenticated resource, **Then** they are denied.

---

### User Story 4 - Reset a forgotten password by email (Priority: P1)

On the sign-in path, a user who does not remember their password can request a reset. They receive an email and set a new password, then sign in with that new password. This is not magic-link sign-in.

**Why this priority**: PRODUCT.md BR-052; users otherwise cannot recover the account.

**Independent Test**: Request reset for a known account; complete reset; old password fails; new password signs in.

**Acceptance Scenarios**:

1. **Given** a signed-up user on the sign-in screen, **When** they choose forgot password and submit their email, **Then** they are told a reset email will be sent if that account exists (same generic message if the email is unknown).
2. **Given** a valid reset from that email, **When** they set a new password, **Then** they can sign in with the new password and cannot sign in with the old one.
3. **Given** an expired or already-used reset, **When** they try to set a password, **Then** the product refuses in plain language and they can request a new reset.

---

### Edge Cases

- Email comparison is case-insensitive for sign-in uniqueness.
- Password MUST NOT be returned in API responses.
- Sign-up MUST NOT grant membership in any existing organization.
- Reset email MUST NOT sign the user in by itself (not a magic link).
- Requesting reset for an unknown email MUST NOT reveal that the account does not exist.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow self-serve sign-up with email, password, name, and free-text address (PRODUCT.md BR-004, BR-023, BR-045).
- **FR-002**: System MUST authenticate with email and password only in this slice (no SSO, SAML, OAuth social, magic link).
- **FR-003**: System MUST reject duplicate emails on sign-up.
- **FR-004**: System MUST keep the user authenticated after sign-in until sign-out or session end.
- **FR-005**: System MUST allow sign-out.
- **FR-006**: Unauthenticated callers MUST NOT access organization or bill data (none of that domain is required in this slice beyond denying it).
- **FR-007**: Sign-up MUST NOT set application admin.
- **FR-008**: Validation failures MUST be explained in plain language.
- **FR-009**: Sign-up and sign-in MUST NOT require email verification (BR-048).
- **FR-010**: Sign-in MUST offer **forgot password**. The user MUST be able to reset the password **by email** and then sign in with the new password (BR-052).
- **FR-011**: Completing a reset MUST invalidate the previous password and MUST NOT leave the user signed in solely because they opened the email.

### Key Entities

- **User account**: email, password credential, name, address, not app admin by default.
- **Password reset**: email-delivered, single-use, time-limited request to set a new password.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A new user can sign up and sign in in one sitting without SSO.
- **SC-002**: Wrong password never yields an authenticated session.
- **SC-003**: Signed-out users cannot call authenticated endpoints successfully.
- **SC-004**: 100% of new sign-ups have name and address stored and isAppAdmin false.
- **SC-005**: A user who forgot the password can complete reset-by-email and sign in with the new password; the old password no longer works.

## Assumptions

- Session/token implementation is chosen in the plan for this feature; product rule is email+password only.
- Minimum password length defaults to 8 characters if the plan does not specify otherwise.
- Reset link expiry defaults to 1 hour if the plan does not specify otherwise.
- Local/dev may capture reset emails (e.g. a mail catcher); production mail is a later deploy concern.
