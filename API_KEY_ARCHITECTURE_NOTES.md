# API Key Management Architecture - Decision Notes

## Context
The NatureDopes API needs a system for distributing API keys to allow programmatic access to flora observation data. The key question is: **Who should be able to generate API keys?**

---

## Approach 1: Open API Key Generation (Current Implementation)

### How It Works
- Anyone can hit `POST /api/keys` and generate an API key
- No authentication required to create, view, or manage keys
- Keys act as the only authentication layer

### Pros
✅ **Simplicity**: No user authentication system needed in Go API
✅ **Friction-free**: Users/researchers can get started immediately
✅ **Decoupled**: API is completely independent from the Next.js app
✅ **Public data mindset**: If flora observations are meant to be public/open data, this aligns with open science principles
✅ **Developer-friendly**: Easy for external developers to integrate
✅ **No single point of failure**: Next.js app going down doesn't prevent API access

### Cons
❌ **No accountability**: Can't track who owns which key
❌ **Abuse potential**: Unlimited key generation could lead to spam/abuse
❌ **No key recovery**: If user loses their key, they can't retrieve it
❌ **Rate limiting challenges**: Can't implement per-user quotas, only per-key
❌ **No revocation by user**: Users can't manage their own keys (would need to contact you)
❌ **Analytics gaps**: Can't analyze usage by user demographic

### Best For
- Open data initiatives
- Public APIs with read-only access
- Community science projects
- When you want minimal barriers to data access

### Security Considerations
- Need aggressive rate limiting per IP and per key
- Should implement key usage quotas (e.g., 1000 requests/day per key)
- Consider CAPTCHA on key generation endpoint to prevent bot abuse
- Log key creation with IP addresses for abuse investigation

---

## Approach 2: Authenticated Users Only

### How It Works
- Users must register/login via Next.js app
- Authenticated users can generate keys through a protected page
- Keys are linked to user accounts (`api_keys.user_id`)
- API key management endpoints require user authentication

### Pros
✅ **Accountability**: Every key is tied to a user
✅ **User self-service**: Users can manage/revoke their own keys
✅ **Better analytics**: Can track usage by user, location, research group, etc.
✅ **Abuse prevention**: Can ban users, limit keys per user
✅ **Relationship building**: You know who's using your API
✅ **Key recovery**: Users can regenerate keys if lost
✅ **Granular permissions**: Could later add user roles (researcher, student, admin)
✅ **Terms of Service**: Can require ToS acceptance during registration

### Cons
❌ **Higher friction**: Users must create account before accessing API
❌ **Coupling**: API depends on user authentication system
❌ **Complexity**: Need to sync user sessions between Next.js and Go
❌ **Maintenance**: Two auth systems to maintain (NextAuth + Go API auth)
❌ **Privacy concerns**: Users must provide email/personal info

### Best For
- APIs with sensitive data
- When you want to build a community/user base
- Research projects requiring attribution
- When quota management per user is important

### Implementation Requirements
1. Database migration: Add `user_id` to `api_keys` table
2. User authentication middleware in Go API
3. Protected Next.js page for key management
4. Session validation between Next.js and Go

---

## Approach 3: Hybrid - Tiered Access

### How It Works
- **Tier 1 (Anonymous)**: Anyone can generate a limited API key
  - Rate limit: 100 requests/day
  - No key management (can't view/revoke)
  - Keys expire after 30 days

- **Tier 2 (Authenticated)**: Registered users get enhanced keys
  - Rate limit: 10,000 requests/day
  - Can manage multiple keys
  - Keys don't expire
  - Can regenerate if lost
  - Usage analytics dashboard

### Pros
✅ **Best of both worlds**: Low friction for casual users, features for power users
✅ **Natural upsell**: Anonymous users see value, then register for more
✅ **Flexibility**: Users choose their level of commitment
✅ **Experimentation**: Researchers can test before committing
✅ **Gradual trust**: Build trust with anonymous tier first

### Cons
❌ **Complexity**: Two systems to build and maintain
❌ **Confusion**: Users might not understand tier differences
❌ **Abuse potential**: Still vulnerable to anonymous tier abuse
❌ **Support burden**: More combinations of issues to troubleshoot

---

## Approach 4: Invite/Request System

### How It Works
- Users submit a request form explaining their use case
- You manually review and approve
- Approved users get API keys via email
- Keys are associated with approved requests

### Pros
✅ **Maximum control**: You vet every user
✅ **Research opportunities**: Learn how people want to use the API
✅ **Network building**: Personal connection with each user
✅ **Quality over quantity**: Fewer, more serious users

### Cons
❌ **Manual work**: You become a bottleneck
❌ **Slow**: Users wait for approval (hours/days)
❌ **Doesn't scale**: Not viable if API becomes popular
❌ **Barrier to entry**: Deters casual/exploratory use

---

## Key Questions to Consider

### 1. **What's the nature of your data?**
- Is flora location data sensitive or public?
- Could malicious actors misuse the data (e.g., overharvesting rare species)?
- Is there privacy risk for users who submitted observations?

### 2. **What's your goal for the API?**
- **Open science/education**: Lean toward open access
- **Research collaboration**: Lean toward authenticated
- **Monetization someday**: Authenticated is easier to add paid tiers to
- **Community building**: Authenticated helps build relationships

### 3. **What's your capacity for abuse management?**
- Can you monitor and respond to abuse quickly?
- Do you have time to implement rate limiting and security measures?
- Are you okay with occasional cleanup work?

### 4. **What's your technical preference?**
- Comfortable with maintaining two auth systems?
- Want to keep Go API simple and stateless?
- Willing to add complexity for better control?

### 5. **What's the user journey you envision?**
```
Journey A (Open):
Researcher hears about API → Visits docs → Generates key → Starts using → Success

Journey B (Authenticated):
Researcher hears about API → Visits site → Sees login requirement →
Creates account → Generates key → Starts using → Success

Journey C (Hybrid):
Researcher hears about API → Generates anonymous key → Tests it →
Likes it → Registers for better limits → Success
```

Which feels right for your community?

---

## Recommendation Framework

### Choose **OPEN** if:
- Flora data is meant to be fully public
- You prioritize ease of access over control
- You're comfortable with aggressive rate limiting as the main defense
- You want the API to "just work" independently
- You're okay not knowing who uses the API

### Choose **AUTHENTICATED** if:
- You want to build relationships with your user base
- You need accountability for data usage
- You plan to add features that require user accounts anyway
- You want detailed usage analytics
- You're comfortable with the added complexity

### Choose **HYBRID** if:
- You want to attract casual users while rewarding serious ones
- You have time to build a more complex system
- You want to test the waters with anonymous tier first
- You like the freemium model approach

### Choose **INVITE-ONLY** if:
- You're in beta/testing phase
- You want to curate your user community
- Data is semi-sensitive
- You have time for manual approval

---

## Middle-Ground Suggestion

If you're unsure, consider this **staged rollout**:

### Phase 1 (Now): Open with Safety Nets
- Keep current open system
- Add aggressive rate limiting (50 requests/hour per IP, 1000/day per key)
- Add CAPTCHA to key generation
- Log all key creation with IP/timestamp
- Set keys to expire after 90 days
- Monitor for abuse

### Phase 2 (Later): Add Authentication Option
- Build the authenticated user key management page
- Authenticated users get higher limits and permanent keys
- Anonymous keys still work but with restrictions
- See which users prefer which method

### Phase 3 (Future): Decide Based on Data
- After 3-6 months, evaluate:
  - How much abuse occurred?
  - Did authenticated tier get adoption?
  - What do users prefer?
- Then commit to a long-term approach

This lets you learn from real usage without committing too early.

---

## Security Implementation Regardless of Approach

**Must-haves for ANY approach:**
1. **Rate limiting**: Per IP and per key
2. **Key expiration**: At least for anonymous keys
3. **Usage logging**: Track requests per key
4. **Abuse detection**: Alert on suspicious patterns
5. **CORS configuration**: Restrict to known domains if possible
6. **Input validation**: Prevent injection attacks
7. **HTTPS only**: No API access over HTTP

**Nice-to-haves:**
- Key rotation reminders
- Usage analytics dashboard
- Webhook for abuse alerts
- IP whitelist option for trusted users
- Read-only scopes (if you add write endpoints later)

---

## Questions to Reflect On

1. **Philosophical**: Is NatureDopes data a public good that should be freely accessible, or a community resource that requires membership?

2. **Practical**: Do you want to know your users, or prefer privacy/anonymity?

3. **Growth**: If 1000 people want API access tomorrow, which system would you rather manage?

4. **Abuse**: If someone generates 10,000 keys and hammers your API, how do you want to handle it?

5. **Future**: Where do you see the API in 2 years? Still free and open, or potentially monetized/restricted?

6. **Community**: Do you want API users to feel like part of the NatureDopes community, or just consumers of data?

---

## My Observation

Based on the project so far:
- You've built a user system in Next.js (registration, login, password reset)
- You're using Prisma for user management
- The app has user-submitted images in the gallery
- This suggests you're building a **community platform**, not just a data dump

If that's true, **authenticated API keys** might align better with your overall vision. BUT if you want the API to be a separate, open-data initiative for science/education, then **open access** makes sense.

The hybrid approach gives you both, but at the cost of complexity.

---

## Next Steps (Once You Decide)

Let me know which direction feels right, and I can:
1. Implement the chosen approach
2. Update the database schema if needed
3. Build the necessary endpoints
4. Create the Next.js page for key management
5. Update the claude.md with the final architecture

Take your time thinking this through - it's a foundational decision that's easier to get right now than to change later!
