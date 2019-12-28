<template>
  <b-container>
    <b-row class="mb-2"><b-col><h2>Profile</h2></b-col></b-row>
    <b-row>
      <b-col cols=4>Issuing authority</b-col>
      <b-col>{{$api.profile.iss}}</b-col>
    </b-row>
    <b-row>
      <b-col cols=4>Audience</b-col>
      <b-col>{{$api.profile.aud}}</b-col>
    </b-row>
    <b-row>
      <b-col cols=4>Subject</b-col>
      <b-col>{{$api.profile.sub}}</b-col>
    </b-row>
    <b-row>
      <b-col cols=4>Name</b-col>
      <b-col>{{$api.profile.name}}</b-col>
    </b-row>
    <b-row>
      <b-col cols=4>Email</b-col>
      <b-col>{{$api.profile.email}}</b-col>
    </b-row>
    <b-row>
      <b-col cols=4>Groups</b-col>
      <b-col>{{$api.profile.groups}}</b-col>
    </b-row>
    <b-row class="mt-2">
      <b-col>
        <b-button size="sm" variant="primary" @click="$api.logout()">Logout</b-button>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
export default {
  data() {
    return {
    };
  },
};
</script>
